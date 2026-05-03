package cluster

import (
	"encoding/json"
	"errors"
	"fmt"
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

	//IpAddressesJson
	if field.IpAddressesJson != "" {
		var ipArray []form.ClusterNodeIpAddr
		if err := json.Unmarshal([]byte(field.IpAddressesJson), &ipArray); err != nil {
			common.ErrorResp(c, errors.New("invalid ip_addresses_json format: "+field.IpAddressesJson), -1)
			return
		}

		fmt.Println("ipArray:", ipArray)

		for _, ipinfo := range ipArray {
			fmt.Println("ipinfo:", ipinfo)

			common_ip_data := &model.ClusterNodeIpaddr{
				NodeID:         field.ID,
				Description:    ipinfo.Description,
				CanAccess:      ipinfo.CanAccess,
				CanHealthCheck: ipinfo.CanHealthCheck,
				IsHealthy:      true,
				IsOn:           ipinfo.IsOn,
				IsUp:           true,
				Order:          1,
				Status:         true,
				IsDeleted:      false,
				UpdateTime:     time.Now().Unix(),
				CreateTime:     time.Now().Unix(),
			}

			fmt.Println(common_ip_data)

		}
	}

	if field.Name == "" {
		common.ErrorResp(c, errors.New("节点名称不能空!"), -1)
		return
	}

	common_data := &model.ClusterNode{
		Name:            field.Name,
		IpAddressesJson: "",
		UpdateTime:      time.Now().Unix(),
	}

	if err := db.GetDb().Model(&model.ClusterNode{}).Where("id = ?", field.ID).Updates(common_data).Error; err != nil {
		common.ErrorResp(c, err, -1)
		return
	}
	common.SuccessResp(c)
}
