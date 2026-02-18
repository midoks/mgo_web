package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	// "mgo/internal/op"
)

func Details(c *gin.Context) {
	data := common.CommonVer(c)
	c.HTML(http.StatusOK, "backend/admin/admin_details.tmpl", data)
}
