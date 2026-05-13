package system

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
)

type StepItem struct {
	Name string
}

func Transfer(c *gin.Context) {

	data := common.CommonVer(c)
	data["submenu"] = GetSysAdvancedSubMenu()

	steps := []StepItem{
		{Name: "开始"},
		{Name: "迁移数据库"},
		{Name: "迁移API节点"},
		{Name: "变更地址"},
		{Name: "迁移管理平台"},
		{Name: "升级节点配置"},
		{Name: "完成"},
	}
	data["steps"] = steps

	c.HTML(http.StatusOK, "backend/system/transfer/index.tmpl", data)
}
