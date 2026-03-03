package cluster

import (
	"net/http"
	// "strconv"
	// "time"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
	// "mgo/internal/db"
	// "mgo/internal/model"
	// utils "mgo/internal/utils"
)

func GetSettingSubMenu() []form.SubSettingMenu {
	menu := []form.SubSettingMenu{
		{
			Number: 1,
			Name:   "基础设置",
			Link:   "clusters/cluster/settings",
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
			Link:   "clusters/cluster/settings",
			Type:   "a",
		},
	}
	return menu
}

func ClusterSettings(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetSubMenu()
	data["setting_menu"] = GetSettingSubMenu()
	data["cluster_id"] = c.Query("cluster_id")
	c.HTML(http.StatusOK, "backend/cluster/settings/index.tmpl", data)
}
