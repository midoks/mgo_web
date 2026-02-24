package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
	"mgo/internal/db"
	"mgo/internal/model"
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

	session := sessions.Default(c)
	uid := session.Get("user_id")
	adminID := common.ParseAdminId(uid)

	common_data := &model.Admin{
		FullName:   field.Name,
		UpdateTime: time.Now(),
	}
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

func LoginLogs(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetSysBaseSubMenu()
	c.HTML(http.StatusOK, "backend/system/settings_login_logs.tmpl", data)
}
