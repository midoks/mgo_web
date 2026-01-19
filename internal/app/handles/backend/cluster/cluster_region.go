package cluster

import (
	"net/http"
	"strconv"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
	"mgo/internal/db"

	"github.com/gin-gonic/gin"
)

func ClusterRegions(c *gin.Context) {
	data := common.CommonVer(c)
	c.HTML(http.StatusOK, "backend/cluster/regions.tmpl", data)
}

func ClusterRegionsAdd(c *gin.Context) {
	data := common.CommonVer(c)
	c.HTML(http.StatusOK, "backend/cluster/regions_add.tmpl", data)
}

func ClusterRegionsNodes(c *gin.Context) {
	data := common.CommonVer(c)
	c.HTML(http.StatusOK, "backend/cluster/regions_nodes.tmpl", data)
}

func PostClusterRegionsNodesAdd(c *gin.Context) {
	var field form.ClusterRegionAdd
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, 0)
		return
	}

	if field.ID != "" {
		id, _ := strconv.ParseInt(field.ID, 10, 64)
		err := db.UpdateClusterRegion(field.Name, field.Mark, id)
		if err == nil {
			common.SuccessResp(c)
			return
		}
		common.ErrorResp(c, err, 0)
		return
	}

	err := db.AddClusterRegion(field.Name, field.Mark)
	if err == nil {
		common.SuccessResp(c)
		return
	}
	common.ErrorResp(c, err, 0)
}
