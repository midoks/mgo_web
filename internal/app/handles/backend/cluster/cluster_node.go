package cluster

import (
	"encoding/json"
	"errors"
	"strconv"

	// "fmt"
	"net/http"
	// "strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
	"mgo/internal/db"
	"mgo/internal/model"
	tools "mgo/internal/utils"
)

func parseClusterNodeIpArray(ipJson string) ([]form.ClusterNodeIpAddr, error) {
	if ipJson == "" {
		return nil, nil
	}
	var ipArray []form.ClusterNodeIpAddr
	if err := json.Unmarshal([]byte(ipJson), &ipArray); err != nil {
		return nil, errors.New("invalid ip_addresses_json format: " + ipJson)
	}
	return ipArray, nil
}

func GetNodeSubMenu() []form.SubMenu {
	menu := []form.SubMenu{
		{
			Number: 1,
			Name:   "节点看板",
			Link:   "clusters/node/boards",
		},
		{
			Number: 2,
			Name:   "节点详情",
			Link:   "clusters/node/details",
		},
		{
			Number: 3,
			Name:   "运行日志",
			Link:   "clusters/node/logs",
		},
		{
			Number: 4,
			Name:   "安装节点",
			Link:   "clusters/node/install",
		},
		{
			Number: 5,
			Name:   "节点设置",
			Link:   "clusters/node/settings",
		},
	}
	return menu
}

func Node(c *gin.Context) {
	data := common.CommonVer(c)
	c.HTML(http.StatusOK, "backend/cluster/node.tmpl", data)
}

func CreateNode(c *gin.Context) {
	method := strings.ToUpper(c.Request.Method)
	if method == "POST" {
		PostCreateNode(c)
		return
	}
	data := common.CommonVer(c)
	data["submenu"] = GetSubMenu()
	data["cluster_id"] = c.Query("cluster_id")
	c.HTML(http.StatusOK, "backend/cluster/node/create.tmpl", data)
}

func NodeBoards(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetNodeSubMenu()
	data["node_id"] = c.Query("node_id")
	data["cluster_id"] = c.Query("cluster_id")
	c.HTML(http.StatusOK, "backend/cluster/node/boards.tmpl", data)
}

func NodeDatail(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetNodeSubMenu()
	data["node_id"] = c.Query("node_id")
	data["cluster_id"] = c.Query("cluster_id")

	node_id := c.Query("node_id")
	node_idint, _ := strconv.ParseInt(node_id, 10, 64)
	node_data, _ := db.GetClusterNodeByID(node_idint)
	data["Data"] = node_data

	c.HTML(http.StatusOK, "backend/cluster/node/details.tmpl", data)
}

func NodeList(c *gin.Context) {
	var field form.ClusterNodeList
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}
	result, count, err := db.GetClusterNodeListByClusterID(field.ClusterID, field.Page.Page, field.Page.Limit)
	if err != nil {
		common.ErrorResp(c, err, -2)
		return
	}
	common.SuccessLayuiResp(c, count, "ok", result)
}

func PostCreateNode(c *gin.Context) {
	var field form.ClusterCreateNode
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	if field.Ip == "" {
		common.ErrorResp(c, errors.New("ip address cannot be empty!"), -1)
		return
	}

	secret := tools.RandString(32)
	unique_id := tools.RandString(32)

	nodeip := &model.ClusterNode{
		Name:        field.Name,
		Ip:          field.Ip,
		ClusterID:   field.ClusterID,
		IsInstalled: false,
		Secret:      secret,
		UniqueID:    unique_id,
		CreateTime:  time.Now().Unix(),
		UpdateTime:  time.Now().Unix(),
	}

	if err := db.GetDb().Create(nodeip).Error; err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	if field.IpAddressesJson != "" {
		ipArray, err := parseClusterNodeIpArray(field.IpAddressesJson)
		if err != nil {
			common.ErrorResp(c, err, -1)
			return
		}

		// 先全部软删除该节点的所有旧 IP 地址
		if err := db.ClusterNodeIpaddrSoftDeleteByNodeID(nodeip.ID); err != nil {
			common.ErrorResp(c, err, -1)
			return
		}

		// 遍历新列表，存在就更新，不存在就创建
		for _, ipinfo := range ipArray {

			existing := db.ExistClusterNodeIpaddrByNodeIDAndIp(nodeip.ID, ipinfo.Ip)
			if existing {
				// 存在则更新（包含已软删除的记录）
				updateData := map[string]interface{}{
					"description":      ipinfo.Description,
					"ip":               ipinfo.Ip,
					"can_access":       ipinfo.CanAccess,
					"can_health_check": ipinfo.CanHealthCheck,
					"is_healthy":       true,
					"is_on":            ipinfo.IsOn,
					"is_up":            true,
					"order":            1,
					"is_deleted":       0,
					"update_time":      time.Now().Unix(),
				}
				if err := db.GetDb().Unscoped().Model(&model.ClusterNodeIpaddr{}).Where("node_id = ? AND ip = ?", nodeip.ID, ipinfo.Ip).Updates(updateData).Error; err != nil {
					common.ErrorResp(c, err, -2)
					return
				}
			} else {
				common_ip_data := &model.ClusterNodeIpaddr{
					NodeID:         nodeip.ID,
					Ip:             ipinfo.Ip,
					Description:    ipinfo.Description,
					CanAccess:      ipinfo.CanAccess,
					CanHealthCheck: ipinfo.CanHealthCheck,
					IsHealthy:      true,
					IsOn:           ipinfo.IsOn,
					IsUp:           true,
					Order:          1,
					IsDeleted:      0,
				}

				// 创建新数据时,先查找是否有已删除的数据
				if delete_id, err := db.GetClusterNodeIpaddrDeletedID(); err == nil {
					// 存在则更新（包含已软删除的记录）
					updateData := map[string]interface{}{
						"node_id":          nodeip.ID,
						"description":      ipinfo.Description,
						"ip":               ipinfo.Ip,
						"can_access":       ipinfo.CanAccess,
						"can_health_check": ipinfo.CanHealthCheck,
						"is_healthy":       true,
						"is_on":            ipinfo.IsOn,
						"is_up":            true,
						"order":            1,
						"is_deleted":       0,
						"update_time":      time.Now().Unix(),
					}
					if err := db.GetDb().Unscoped().Model(&model.ClusterNodeIpaddr{}).Where("id", delete_id).Updates(updateData).Error; err != nil {
						common.ErrorResp(c, err, -2)
						return
					}

				} else {

					// 不存在则创建
					common_ip_data.CreateTime = time.Now().Unix()
					common_ip_data.UpdateTime = time.Now().Unix()
					if err := db.GetDb().Create(common_ip_data).Error; err != nil {
						common.ErrorResp(c, err, -1)
						return
					}
				}
			}
		}
	}

	common.SuccessResp(c)
}

func PostDeleteNode(c *gin.Context) {
	var field form.ID
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	err := db.ClusterNodeDeleteByID(nil, field.ID)
	if err == nil {
		common.SuccessResp(c)
		return
	}
	common.ErrorResp(c, err, -1)

}
