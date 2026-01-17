package cluster

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
	"mgo/internal/db"
)

func ClusterGroups(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetSubMenu()
	data["cluster_id"] = c.Query("cluster_id")
	c.HTML(http.StatusOK, "backend/cluster/cluster_groups.tmpl", data)
}

func ClusterGroupsAdd(c *gin.Context) {
	data := common.CommonVer(c)
	data["id"] = c.Query("id")
	data["submenu"] = GetSubMenu()
	data["cluster_id"] = c.Query("cluster_id")
	c.HTML(http.StatusOK, "backend/cluster/cluster_groups_add.tmpl", data)
}

func ClusterGroupsList(c *gin.Context) {
	result, count, _ := db.GetClusterGroupList(1, 10)
	common.SuccessLayuiResp(c, count, "ok", result)
}

func PostClusterGroupsAdd(c *gin.Context) {
	var field form.ClusterGroupAdd
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, 0)
		return
	}

	err := db.AddClusterGroup(field.Name, field.ClusterId)
	if err == nil {
		common.SuccessResp(c)
		return
	}
	common.ErrorResp(c, err, 0)
}

func ClusterGroupsDelete(c *gin.Context) {
	var field form.ID
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	fmt.Println(field.ID)
	err := db.ClusterGroupDeleteById(field.ID)
	if err == nil {
		common.SuccessResp(c)
		return
	}
	common.ErrorResp(c, err, -1)
}
