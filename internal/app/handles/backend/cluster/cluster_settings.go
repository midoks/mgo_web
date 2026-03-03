package cluster

import (
	"net/http"
	"strconv"
	// "time"
	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
	"mgo/internal/db"
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
			Name:   "健康检查",
			Link:   "clusters/cluster/settings/health",
			Type:   "a",
		},
	}
	return menu
}

func ClusterSettings(c *gin.Context) {
	cluster_id := c.Query("cluster_id")
	data := common.CommonVer(c)
	data["submenu"] = GetSubMenu()
	data["setting_menu"] = GetSettingSubMenu()
	data["cluster_id"] = cluster_id

	cluster_idint, _ := strconv.ParseInt(cluster_id, 10, 64)
	cluster_data, _ := db.GetClusterByID(cluster_idint)

	data["Data"] = cluster_data
	c.HTML(http.StatusOK, "backend/cluster/settings/index.tmpl", data)
}
