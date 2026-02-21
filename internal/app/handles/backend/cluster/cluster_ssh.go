package cluster

import (
	"net/http"
	// "strconv"
	"time"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
	"mgo/internal/db"
	"mgo/internal/model"
)

func GetSshSubMenu() []form.ClusterSubMenu {
	menu := []form.ClusterSubMenu{
		{
			Number: 1,
			Name:   "SSH认证列表",
			Link:   "clusters/ssh",
		},
		{
			Number: 2,
			Name:   "创建认证",
			Link:   "clusters/ssh/create",
		},
	}
	return menu
}

func ClusterSsh(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetSshSubMenu()
	c.HTML(http.StatusOK, "backend/cluster/ssh.tmpl", data)
}

func ClusterSshCreate(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetSshSubMenu()
	c.HTML(http.StatusOK, "backend/cluster/ssh_create.tmpl", data)
}

func ClusterSshList(c *gin.Context) {
	var field form.Page
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	result, count, _ := db.GetClusterSshList(field.Page, field.Limit)
	common.SuccessLayuiResp(c, count, "ok", result)
}

func PostClusterSshCreate(c *gin.Context) {
	var field form.ClusterSshCreate
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	common_data := &model.ClusterSsh{
		Name:           field.Name,
		Method:         field.Method,
		Username:       field.Username,
		Password:       field.Password,
		Privatekey:     field.Privatekey,
		PrivatekeyPass: field.PrivatekeyPass,
		Mark:           field.Mark,
		UpdateTime:     time.Now(),
	}

	if field.ID > 0 {
		if err := db.GetDb().Model(&model.ClusterSsh{}).Where("id = ?", field.ID).Updates(common_data).Error; err != nil {
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
