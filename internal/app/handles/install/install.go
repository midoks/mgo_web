package install

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	// "mgo/internal/conf"
)

func Home(c *gin.Context) {
	data := common.CommonVer()
	c.HTML(http.StatusOK, "install/index", data)
}
