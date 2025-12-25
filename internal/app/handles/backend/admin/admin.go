package admin

import (
	"fmt"
	"net/http"
	"strconv"

	// "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/db"
	// "mgo/internal/op"
)

func Home(c *gin.Context) {
	data := common.CommonVer()
	c.HTML(http.StatusOK, "backend/admin/index.tmpl", data)
}

func Edit(c *gin.Context) {

	id := c.Query("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)

	fmt.Println("id:", id)

	admin_data, _ := db.GetAdminById(idInt)
	fmt.Println("admin_data:", admin_data)

	data := common.CommonVer()

	data["Data"] = admin_data
	c.HTML(http.StatusOK, "backend/admin/edit.tmpl", data)
}

func List(c *gin.Context) {
	result, count, _ := db.GetAdminList(1, 10)
	common.SuccessLayuiResp(c, count, "ok", result)
}
