package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	// "mgo/internal/op"
)

func Advanced(c *gin.Context) {
	data := common.CommonVer(c)
	c.HTML(http.StatusOK, "backend/system/advanced.tmpl", data)
}
