package admin

import (
	// "fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
	"mgo/internal/db"
	"mgo/internal/model"
)

// 通知媒介
func Recipients(c *gin.Context) {
	data := common.CommonVer(c)
	c.HTML(http.StatusOK, "backend/admin/recipients.tmpl", data)
}

func RecipientsInstances(c *gin.Context) {
	data := common.CommonVer(c)
	c.HTML(http.StatusOK, "backend/admin/recipients_instances.tmpl", data)
}

func RecipientsInstancesAdd(c *gin.Context) {
	data := common.CommonVer(c)
	c.HTML(http.StatusOK, "backend/admin/recipients_instances_add.tmpl", data)
}

func RecipientsInstancesDetails(c *gin.Context) {
	data := common.CommonVer(c)
	data["id"] = c.Query("id")
	c.HTML(http.StatusOK, "backend/admin/recipients_instances_details.tmpl", data)
}

func RecipientsInstancesUpdate(c *gin.Context) {
	data := common.CommonVer(c)
	data["id"] = c.Query("id")
	c.HTML(http.StatusOK, "backend/admin/recipients_instances_update.tmpl", data)
}

func RecipientsInstancesTest(c *gin.Context) {
	data := common.CommonVer(c)
	data["id"] = c.Query("id")
	c.HTML(http.StatusOK, "backend/admin/recipients_instances_test.tmpl", data)
}

func RecipientsList(c *gin.Context) {
	result, count, _ := db.GetAdminRecipientsList(1, 10)
	common.SuccessLayuiResp(c, count, "ok", result)
}

func PostRecipientsInstancesAdd(c *gin.Context) {
	var field form.AdminRecipients
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	admin_recipicents := &model.AdminMediaInstance{
		Name:       field.Name,
		MediaType:  field.MediaType,
		State:      true,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	if err := db.GetDb().Create(admin_recipicents).Error; err != nil {
		common.ErrorResp(c, err, -1)
		return
	}
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
