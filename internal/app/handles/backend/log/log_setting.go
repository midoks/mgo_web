package log

import (
	"fmt"
	// "strconv"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
	"mgo/internal/db"
	"mgo/internal/model"
)

const LOGSYSNAME = "log_sys"

func Settings(c *gin.Context) {

	log_data, _ := db.GetSysSettingByCode(LOGSYSNAME)
	// fmt.Println(log_data)
	data := common.CommonVer(c)
	data["submenu"] = GetLogSubMenu()
	data["Data"] = log_data
	c.HTML(http.StatusOK, "backend/log/setting.tmpl", data)
}

func PostSettting(c *gin.Context) {
	var field form.LogSetting
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, 0)
		return
	}

	fmt.Println(field)
	fmt.Println("field.AllowedManualDelete:", field.AllowedManualDelete, "field.AllowedManual:", field.AllowedManual, "field.SaveDay:", field.SaveDay)

	common_data := &model.SysSetting{
		Code:       LOGSYSNAME,
		UpdateTime: time.Now(),
	}

	common_data.Uid = 0
	common_data.SetLogValue(model.SysSettingLogValue{
		AllowedManualDelete: field.AllowedManualDelete,
		AllowedManual:       field.AllowedManual,
		SaveDay:             field.SaveDay,
		MaxCapacityLimit:    field.MaxCapacityLimit,
		MaxCapacityUnit:     field.MaxCapacityUnit,
	})

	_, err := db.GetSysSettingByCode(LOGSYSNAME)
	if err == nil {
		if err := db.GetDb().Model(&model.SysSetting{}).Where("code = ?", LOGSYSNAME).Updates(common_data).Error; err != nil {
			common.ErrorResp(c, err, -1)
			return
		}
		common.SuccessResp(c)
		return
	}

	common_data.CreateTime = time.Now()
	if err := db.GetDb().Create(common_data).Error; err != nil {
		common.ErrorResp(c, err, -1)
		return
	}
	common.SuccessResp(c)
}
