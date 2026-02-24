package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
	"mgo/internal/db"
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

func Login(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetSysBaseSubMenu()
	c.HTML(http.StatusOK, "backend/system/settings_login.tmpl", data)
}

func List(c *gin.Context) {
	result, count, _ := db.GetAdminList(1, 10)
	common.SuccessLayuiResp(c, count, "ok", result)
}
