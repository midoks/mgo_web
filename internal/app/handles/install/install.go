package install

import (
	// "fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	// "mgo/internal/conf"
)

func HomePage(c *gin.Context) {
	data := common.CommonVer()
	step := c.Query("step")
	if step == "2" {
		c.HTML(http.StatusOK, "install/step2.tmpl", data)
		return
	}
	c.HTML(http.StatusOK, "install/index.tmpl", data)
}

func Step2Page(c *gin.Context) {
	data := common.CommonVer()
	c.HTML(http.StatusOK, "install/step2.tmpl", data)
}
