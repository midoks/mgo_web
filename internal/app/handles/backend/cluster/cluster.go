package cluster

import (
	"net/http"
	// "strconv"
	"time"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
	"mgo/internal/db"
	"mgo/internal/model"
	tools "mgo/internal/utils"
)

func GetSubMenu() []form.SubMenu {
	menu := []form.SubMenu{
		{
			Number: 1,
			Name:   "集群看板",
			Link:   "clusters/cluster/boards",
			Icon:   "layui-icon-home",
		},
		{
			Number: 2,
			Name:   "节点列表",
			Link:   "clusters/cluster/list",
			Icon:   "layui-icon-list",
		},
		{
			Number: 3,
			Name:   "创建节点",
			Link:   "clusters/cluster/create_node",
			Icon:   "layui-icon-add-circle",
		},
		{
			Number: 4,
			Name:   "安装升级",
			Link:   "clusters/cluster/upgrade",
			Icon:   "layui-icon-download-circle",
		},
		{
			Number: 5,
			Name:   "节点分组",
			Link:   "clusters/cluster/groups",
			Icon:   "layui-icon-group",
		},
		{
			Number: 6,
			Name:   "集群设置",
			Link:   "clusters/cluster/settings",
			Icon:   "layui-icon-set",
		},
		{
			Number: 7,
			Name:   "其它操作",
			Link:   "clusters/cluster/delete",
			Icon:   "layui-icon-more-vertical",
		},
	}
	return menu
}

func Home(c *gin.Context) {
	data := common.CommonVer(c)
	c.HTML(http.StatusOK, "backend/cluster/index.tmpl", data)
}

func Create(c *gin.Context) {
	data := common.CommonVer(c)
	c.HTML(http.StatusOK, "backend/cluster/create.tmpl", data)
}

func ClusterBoards(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetSubMenu()
	data["cluster_id"] = c.Query("cluster_id")
	c.HTML(http.StatusOK, "backend/cluster/boards.tmpl", data)
}

func ClusterList(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetSubMenu()
	data["cluster_id"] = c.Query("cluster_id")
	c.HTML(http.StatusOK, "backend/cluster/list.tmpl", data)
}

func ClusterInstall(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetSubMenu()
	data["cluster_id"] = c.Query("cluster_id")
	c.HTML(http.StatusOK, "backend/cluster/install.tmpl", data)
}

func ClusterUpgrade(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetSubMenu()
	data["cluster_id"] = c.Query("cluster_id")
	c.HTML(http.StatusOK, "backend/cluster/install_upgrade.tmpl", data)
}

func ClusterDelete(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetSubMenu()
	data["cluster_id"] = c.Query("cluster_id")
	c.HTML(http.StatusOK, "backend/cluster/delete.tmpl", data)
}

func List(c *gin.Context) {
	var field form.Page
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	result, count, _ := db.GetClusterList(field.Page, field.Limit)
	common.SuccessLayuiResp(c, count, "ok", result)
}

func PostCreate(c *gin.Context) {
	var field form.ClusterCreate
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	unique_id := tools.RandString(32)
	cluster := &model.Cluster{
		Name:       field.Name,
		UniqueID:   unique_id,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}

	if err := db.GetDb().Create(cluster).Error; err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	common.SuccessResp(c)
}

func Delete(c *gin.Context) {
	var field form.ID
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	err := db.ClusterDeleteByID(nil, field.ID)
	if err == nil {
		common.SuccessResp(c)
		return
	}
	common.ErrorResp(c, err, -1)
}
