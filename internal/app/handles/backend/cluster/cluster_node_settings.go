package cluster

import (
	// "errors"
	// "fmt"
	"net/http"
	"strconv"

	// "strings"
	// "time"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
	"mgo/internal/db"
	// "mgo/internal/model"
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
	c.HTML(http.StatusOK, "backend/cluster/node/settings.tmpl", data)
}

func NodeSettingsSsh(c *gin.Context) {
	node_id := c.Query("node_id")
	data := common.CommonVer(c)
	data["submenu"] = GetNodeSubMenu()
	data["setting_menu"] = GetNodeSettingSubMenu()
	data["node_id"] = node_id
	data["cluster_id"] = c.Query("cluster_id")

	node_idint, _ := strconv.ParseInt(node_id, 10, 64)
	node_login_data, err := db.GetClusterNodeLoginByNodeID(node_idint)
	if err != nil {
		c.HTML(http.StatusOK, "backend/cluster/node/settings_ssh.tmpl", data)
		return
	}
	data["Data"] = node_login_data
	node_login_param, err := node_login_data.GetParams()
	if err == nil {
		if node_login_param.SshID > 0 {
			ssh_data, _ := db.GetClusterSshByID(node_login_param.SshID)
			data["SshData"] = ssh_data
		}
	}
	c.HTML(http.StatusOK, "backend/cluster/node/settings_ssh.tmpl", data)
}
