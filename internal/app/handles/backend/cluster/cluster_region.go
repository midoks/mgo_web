package cluster

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
)

func ClusterRegions(c *gin.Context) {
	data := common.CommonVer(c)
	c.HTML(http.StatusOK, "backend/cluster/regions.tmpl", data)
}
