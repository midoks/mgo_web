package cluster

import (
	// "fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/db"
)

func Create(c *gin.Context) {
	data := common.CommonVer(c)
	c.HTML(http.StatusOK, "backend/cluster/create.tmpl", data)
}

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

func ClusterCreateNode(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetSubMenu()
	data["cluster_id"] = c.Query("cluster_id")
	c.HTML(http.StatusOK, "backend/cluster/cluster_create_node.tmpl", data)
}

func Node(c *gin.Context) {
	data := common.CommonVer(c)
	c.HTML(http.StatusOK, "backend/cluster/node.tmpl", data)
}
