package admin

import (
	// "fmt"
	"net/http"

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

func List(c *gin.Context) {
	result, count, _ := db.GetAdminList(1, 10)
	common.SuccessLayuiResp(c, count, "ok", result)
}
