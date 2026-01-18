package cluster

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
)

func ClusterGrants(c *gin.Context) {
	data := common.CommonVer(c)
	c.HTML(http.StatusOK, "backend/cluster/grants.tmpl", data)
}
