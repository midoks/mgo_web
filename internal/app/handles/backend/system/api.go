package system

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/db"
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

func ApiDetails(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetSysAdvancedSubMenu()

	id := c.Query("id")
	data["id"] = id

	idInt, _ := strconv.ParseInt(id, 10, 64)
	dbnode_data, err := db.GetApiNodeByID(idInt)
	if err == nil {
		data["Data"] = dbnode_data
	}

	c.HTML(http.StatusOK, "backend/system/api/details.tmpl", data)
}
