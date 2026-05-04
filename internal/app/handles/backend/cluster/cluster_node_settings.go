package cluster

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
	"mgo/internal/db"
	"mgo/internal/model"
)

func GetNodeSettingSubMenu() []form.SubSettingMenu {
	menu := []form.SubSettingMenu{
		{
			Number: 1,
			Name:   "基础设置",
			Link:   "clusters/node/settings",
			Type:   "a",
		},
		{
			Number: 2,
			Name:   "line",
			Link:   "",
			Type:   "line",
		},
		{
			Number: 3,
			Name:   "SSH设置",
			Link:   "clusters/node/settings/ssh",
			Type:   "a",
		},
	}
	return menu
}

func NodeSettings(c *gin.Context) {
	node_id := c.Query("node_id")

	data := common.CommonVer(c)
	data["submenu"] = GetNodeSubMenu()
	data["setting_menu"] = GetNodeSettingSubMenu()
	data["node_id"] = node_id
	data["cluster_id"] = c.Query("cluster_id")

	node_idint, _ := strconv.ParseInt(node_id, 10, 64)
	node_data, _ := db.GetClusterNodeByID(node_idint)
	data["Data"] = node_data

	ipaddrs, err := db.GetClusterNodeIpaddrByNodeID(node_idint)
	data["IpAddressesJson"] = "[]"
	if err == nil {
		ipaddrs_json, err := json.Marshal(ipaddrs)
		if err == nil {
			data["IpAddressesJson"] = string(ipaddrs_json)
		}
	}

	c.HTML(http.StatusOK, "backend/cluster/node/settings.tmpl", data)
}

func PostNodeSettings(c *gin.Context) {
	var field form.ClusterNodeSettings
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	// IpAddressesJson
	var ipArray []form.ClusterNodeIpAddr
	if field.IpAddressesJson != "" {
		var err error
		ipArray, err = parseClusterNodeIpArray(field.IpAddressesJson)
		if err != nil {
			common.ErrorResp(c, err, -1)
			return
		}

		// 先全部软删除该节点的所有旧 IP 地址
		if err := db.ClusterNodeIpaddrSoftDeleteByNodeID(field.ID); err != nil {
			common.ErrorResp(c, err, -1)
			return
		}

		// 遍历新列表，存在就更新，不存在就创建
		for _, ipinfo := range ipArray {

			existing := db.ExistClusterNodeIpaddrByNodeIDAndIp(field.ID, ipinfo.Ip)
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
				if err := db.GetDb().Unscoped().Model(&model.ClusterNodeIpaddr{}).Where("node_id = ? AND ip = ?", field.ID, ipinfo.Ip).Updates(updateData).Error; err != nil {
					common.ErrorResp(c, err, -2)
					return
				}
			} else {
				common_ip_data := &model.ClusterNodeIpaddr{
					NodeID:         field.ID,
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
						"node_id":          field.ID,
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

	if field.Name == "" {
		common.ErrorResp(c, errors.New("节点名称不能空!"), -1)
		return
	}

	common_data := &model.ClusterNode{
		Name:            field.Name,
		IpAddressesJson: field.IpAddressesJson,
		UpdateTime:      time.Now().Unix(),
	}

	if err := db.GetDb().Model(&model.ClusterNode{}).Where("id = ?", field.ID).Updates(common_data).Error; err != nil {
		common.ErrorResp(c, err, -1)
		return
	}
	common.SuccessResp(c)
}
