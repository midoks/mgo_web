package log

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	// "mgo/internal/app/form"
	// "mgo/internal/db"
)

func Settings(c *gin.Context) {
	data := common.CommonVer(c)

	c.HTML(http.StatusOK, "backend/log/setting.tmpl", data)
}
