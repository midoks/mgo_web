package cluster

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
	"mgo/internal/db"
	"mgo/internal/model"
)

func SelectIp(c *gin.Context) {
	data := common.CommonVer(c)
	region_list, _, _ := db.GetClusterRegionList(1, 100)
	data["region_list"] = region_list
	c.HTML(http.StatusOK, "backend/cluster/cluster_select_ip.tmpl", data)
}

func SelectRegion(c *gin.Context) {
	data := common.CommonVer(c)
	region_list, _, _ := db.GetClusterRegionList(1, 100)
	data["region_list"] = region_list
	c.HTML(http.StatusOK, "backend/cluster/cluster_select_region.tmpl", data)
}

func SelectGroups(c *gin.Context) {
	data := common.CommonVer(c)
	region_list, _, _ := db.GetClusterGroupList(1, 100)
	data["groups_list"] = region_list
	c.HTML(http.StatusOK, "backend/cluster/cluster_select_groups.tmpl", data)
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
	c.HTML(http.StatusOK, "backend/cluster/cluster_create_node.tmpl", data)
}

func Node(c *gin.Context) {
	data := common.CommonVer(c)
	c.HTML(http.StatusOK, "backend/cluster/node.tmpl", data)
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
	data := common.CommonVer(c)
	var field form.ClusterCreateNode
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	if field.Ip == "" {
		common.ErrorResp(c, errors.New("IP不能为空!"), -1)
		return
	}

	fmt.Println("field2:", field.Ip)

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

	// common.SuccessResp(c)
	data["submenu"] = GetSubMenu()
	data["cluster_id"] = field.ClusterID
	c.HTML(http.StatusOK, "backend/cluster/cluster_create_node2.tmpl", data)
}

func PostDeleteNode(c *gin.Context) {
	var field form.ID
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	err := db.ClusterNodeDeleteById(field.ID)
	if err == nil {
		common.SuccessResp(c)
		return
	}
	common.ErrorResp(c, err, -1)

}
