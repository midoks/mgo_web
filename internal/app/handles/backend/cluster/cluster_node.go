package cluster

import (
	"errors"
	// "fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
	"mgo/internal/db"
	"mgo/internal/model"
)

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
	c.HTML(http.StatusOK, "backend/cluster/node/detail.tmpl", data)
}

func NodeInstall(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetNodeSubMenu()
	data["node_id"] = c.Query("node_id")
	data["cluster_id"] = c.Query("cluster_id")
	c.HTML(http.StatusOK, "backend/cluster/node/install.tmpl", data)
}

func NodeLogs(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetNodeSubMenu()
	data["node_id"] = c.Query("node_id")
	data["cluster_id"] = c.Query("cluster_id")
	c.HTML(http.StatusOK, "backend/cluster/node/logs.tmpl", data)
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

	nodeip := &model.ClusterNode{
		Name:       field.Name,
		Ip:         field.Ip,
		ClusterID:  field.ClusterID,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	if err := db.GetDb().Create(nodeip).Error; err != nil {
		common.ErrorResp(c, err, -1)
		return
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
