package log

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
	// "mgo/internal/db"
)

func Settings(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetLogSubMenu()
	c.HTML(http.StatusOK, "backend/log/setting.tmpl", data)
}

func PostSettting(c *gin.Context) {
	var field form.LogSetting
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, 0)
		return
	}

	fmt.Println(field)
	common.SuccessResp(c)
}
