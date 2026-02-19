package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
	// "mgo/internal/op"
)

func GetSysAdvancedSubMenu() []form.ClusterSubMenu {
	menu := []form.ClusterSubMenu{
		{
			Number: 1,
			Name:   "数据库",
			Link:   "clusters/cluster/boards",
		},
		{
			Number: 2,
			Name:   "API节点",
			Link:   "clusters/cluster/list",
		},
		{
			Number: 3,
			Name:   "用户节点",
			Link:   "clusters/cluster/create_node",
		},
		{
			Number: 4,
			Name:   "日志数据库",
			Link:   "clusters/cluster/create_node",
		},
		{
			Number: 5,
			Name:   "迁移",
			Link:   "clusters/cluster/create_node",
		},
	}
	return menu
}

func Advanced(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetSysAdvancedSubMenu()
	c.HTML(http.StatusOK, "backend/system/advanced.tmpl", data)
}
