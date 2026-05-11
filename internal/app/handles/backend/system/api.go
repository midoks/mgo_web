package system

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
)

func Api(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetSysAdvancedSubMenu()
	c.HTML(http.StatusOK, "backend/system/api/index.tmpl", data)
}

func ApiAdd(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetSysAdvancedSubMenu()
	c.HTML(http.StatusOK, "backend/system/api/add.tmpl", data)
}
