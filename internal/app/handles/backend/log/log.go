package log

import (
	// "errors"
	"net/http"
	// "strconv"
	// "time"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
	"mgo/internal/db"
	// "mgo/internal/model"
	// utils "mgo/internal/utils"
	// "mgo/internal/op"
)

func Home(c *gin.Context) {
	data := common.CommonVer(c)

	c.HTML(http.StatusOK, "backend/log/index.tmpl", data)
}

func List(c *gin.Context) {
	var field form.Page
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	result, count, _ := db.GetLogList(field.Page, field.Limit)
	common.SuccessLayuiResp(c, count, "ok", result)
}

func Delete(c *gin.Context) {
	var field form.ID
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	err := db.LogDeleteById(field.ID)
	if err == nil {
		common.SuccessResp(c)
		return
	}
	common.ErrorResp(c, err, -1)
}
