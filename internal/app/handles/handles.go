package handles

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/conf"
)

func AdminPage(c *gin.Context) {
	if !conf.Security.InstallLock {
		c.Redirect(302, "/install")
	}
	data := common.CommonVer()
	c.HTML(http.StatusOK, "index.tmpl", data)
}

func Home(c *gin.Context) {
	// if !conf.Security.InstallLock {
	// 	c.Redirect(302, "/install")
	// }
	data := common.CommonVer()
	c.HTML(http.StatusOK, "backend/index/index.tmpl", data)
}
