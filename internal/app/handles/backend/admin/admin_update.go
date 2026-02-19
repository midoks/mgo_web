package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/db"
	// "mgo/internal/op"
)

func Update(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)
	admin_data, _ := db.GetAdminById(idInt)

	data := common.CommonVer(c)
	data["id"] = id
	data["Data"] = admin_data
	c.HTML(http.StatusOK, "backend/admin/admin_update.tmpl", data)
}
