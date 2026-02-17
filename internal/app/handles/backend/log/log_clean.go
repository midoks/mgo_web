package log

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	// "mgo/internal/app/form"
	// "mgo/internal/db"
)

func Clean(c *gin.Context) {
	data := common.CommonVer(c)

	c.HTML(http.StatusOK, "backend/log/clean.tmpl", data)
}
