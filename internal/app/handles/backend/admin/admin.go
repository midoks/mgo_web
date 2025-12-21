package admin

import (
	// "fmt"
	"net/http"

	// "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	// "mgo/internal/db"
	// "mgo/internal/op"
)

func HomePage(c *gin.Context) {
	data := common.CommonVer()
	c.HTML(http.StatusOK, "backend/admin/index.tmpl", data)
}
