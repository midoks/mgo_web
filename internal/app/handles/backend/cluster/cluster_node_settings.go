package cluster

import (
	// "errors"
	// "fmt"
	"net/http"
	// "strings"
	// "time"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
	// "mgo/internal/db"
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
	data := common.CommonVer(c)
	data["submenu"] = GetNodeSubMenu()
	data["setting_menu"] = GetNodeSettingSubMenu()
	data["node_id"] = c.Query("node_id")
	c.HTML(http.StatusOK, "backend/cluster/node/settings.tmpl", data)
}

func NodeSettingsSsh(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetNodeSubMenu()
	data["setting_menu"] = GetNodeSettingSubMenu()
	data["node_id"] = c.Query("node_id")
	c.HTML(http.StatusOK, "backend/cluster/node/settings_ssh.tmpl", data)
}
