package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
	"mgo/internal/db"
	"mgo/internal/model"
	utils "mgo/internal/utils"
)

func GetSysBaseSubMenu() []form.ClusterSubMenu {
	menu := []form.ClusterSubMenu{
		{
			Number: 1,
			Name:   "管理员界面设置",
			Link:   "system/settings",
		},
		{
			Number: 2,
			Name:   "个人资料",
			Link:   "system/settings/profile",
		},
		{
			Number: 3,
			Name:   "登录设置",
			Link:   "system/settings/login",
		},
	}
	return menu
}

func Home(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetSysBaseSubMenu()
	c.HTML(http.StatusOK, "backend/system/settings.tmpl", data)
}

func Profile(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetSysBaseSubMenu()
	c.HTML(http.StatusOK, "backend/system/settings_profile.tmpl", data)
}

func PostProfile(c *gin.Context) {
	var field form.SettingProfile
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	if field.Name == "" {
		common.ErrorResp(c, errors.New("你的姓名,不能为空!"), -2)
		return
	}

	common_data := &model.Admin{
		FullName:   field.Name,
		UpdateTime: time.Now().Unix(),
	}
	data := common.CommonVer(c)
	adminID := data["login_uid"]
	if err := db.GetDb().Model(&model.Admin{}).Where("id = ?", adminID).Updates(common_data).Error; err != nil {
		common.ErrorResp(c, err, -1)
		return
	}
	common.SuccessResp(c)
}

func Login(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetSysBaseSubMenu()
	c.HTML(http.StatusOK, "backend/system/settings_login.tmpl", data)
}

func PostLogin(c *gin.Context) {
	var field form.SettingLogin
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	if field.Name == "" {
		common.ErrorResp(c, errors.New("你的姓名,不能为空!"), -2)
		return
	}
	common_data := &model.Admin{
		FullName:   field.Name,
		UpdateTime: time.Now().Unix(),
	}

	if field.Password != "" || field.Password2 != "" {
		if field.Password != field.Password2 {
			common.ErrorResp(c, errors.New("两次密码不一致!"), -2)
			return
		}

		salt := utils.RandString(16)
		common_data.Salt = salt
		common_data.Password = model.TwoHashPwd(field.Password, salt)
	}

	data := common.CommonVer(c)
	adminID := data["login_uid"]
	if err := db.GetDb().Model(&model.Admin{}).Where("id = ?", adminID).Updates(common_data).Error; err != nil {
		common.ErrorResp(c, err, -1)
		return
	}
	common.SuccessResp(c)
}

func LoginLogs(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetSysBaseSubMenu()
	c.HTML(http.StatusOK, "backend/system/settings_login_logs.tmpl", data)
}

func LoginLogsList(c *gin.Context) {
	var field form.AdminPage
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	result, count, _ := db.GetAdminLogsListByAdminId(field.AdminId, field.Page.Page, field.Page.Limit)
	common.SuccessLayuiResp(c, count, "ok", result)
}
