package admin

import (
	// "encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
	"mgo/internal/db"
	"mgo/internal/model"
	// "mgo/internal/op"
)

// 通知媒介 -- S
func Recipients(c *gin.Context) {
	data := common.CommonVer(c)
	c.HTML(http.StatusOK, "backend/admin/recipients.tmpl", data)
}

func RecipientsGroups(c *gin.Context) {
	data := common.CommonVer(c)
	c.HTML(http.StatusOK, "backend/admin/recipients_groups.tmpl", data)
}

func RecipientsInstances(c *gin.Context) {
	data := common.CommonVer(c)
	c.HTML(http.StatusOK, "backend/admin/recipients_instances.tmpl", data)
}

func RecipientsInstancesAdd(c *gin.Context) {
	data := common.CommonVer(c)
	c.HTML(http.StatusOK, "backend/admin/recipients_instances_add.tmpl", data)
}

func RecipientsList(c *gin.Context) {
	result, count, _ := db.GetAdminRecipientsList(1, 10)
	common.SuccessLayuiResp(c, count, "ok", result)
}

func PostRecipientsInstancesAdd(c *gin.Context) {
	var field form.AdminRecipients
	if err := c.ShouldBind(&field); err != nil {
		fmt.Println("ccccc")
		common.ErrorResp(c, err, -1)
		return
	}

	fmt.Println("123123")

	admin_recipicents := &model.AdminMediaInstance{
		Name:       field.Name,
		MediaType:  field.MediaType,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	if err := db.GetDb().Create(admin_recipicents).Error; err != nil {
		fmt.Println("err:", err)
		common.ErrorResp(c, err, -1)
		return
	}

	fmt.Println("field:", field)
	common.SuccessResp(c)
}

func RecipientsDelete(c *gin.Context) {
	var field form.ID
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	err := db.AdminRecipientsDeleteById(field.ID)
	if err == nil {
		common.SuccessResp(c)
		return
	}
	common.ErrorResp(c, err, -1)
}
