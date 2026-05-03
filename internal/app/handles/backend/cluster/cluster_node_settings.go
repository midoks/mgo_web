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

	if node_data != nil && node_data.IpAddressesJson != "" {
		data["IpAddressesJson"] = node_data.IpAddressesJson
	} else {
		data["IpAddressesJson"] = "[]"
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
	if field.IpAddressesJson != "" {
		var ipArray []form.ClusterNodeIpAddr
		if err := json.Unmarshal([]byte(field.IpAddressesJson), &ipArray); err != nil {
			common.ErrorResp(c, errors.New("invalid ip_addresses_json format: "+field.IpAddressesJson), -1)
			return
		}

		// 先全部软删除该节点的所有旧 IP 地址
		if err := db.GetDb().Model(&model.ClusterNodeIpaddr{}).Where("node_id = ?", field.ID).Updates(map[string]interface{}{
			"is_deleted":  1,
			"update_time": time.Now().Unix(),
		}).Error; err != nil {
			common.ErrorResp(c, err, -1)
			return
		}

		// 遍历新列表，存在就更新，不存在就创建
		for _, ipinfo := range ipArray {

			// 查询是否存在（包括已软删除的）
			var existing model.ClusterNodeIpaddr
			err := db.GetDb().Unscoped().Where("node_id = ? AND ip = ?", field.ID, ipinfo.Ip).First(&existing).Error

			if err == nil {
				// 存在则更新（包含已软删除的记录）
				updateData := map[string]interface{}{
					"description":      ipinfo.Description,
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
