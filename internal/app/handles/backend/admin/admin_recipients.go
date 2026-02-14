package admin

import (
	// "fmt"
	"net/http"
	"strconv"
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
	id := c.Query("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)
	recipient_data, _ := db.GetAdminRecipientById(idInt)

	data := common.CommonVer(c)
	data["id"] = id
	data["Data"] = recipient_data
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

	if field.ID > 0 {
		update_data := &model.AdminMediaInstance{
			Name:       field.Name,
			Status:     field.Status,
			Mark:       field.Mark,
			HashLife:   field.HashLife,
			UpdateTime: time.Now(),
		}

		if err := db.GetDb().Model(&model.AdminMediaInstance{}).Where("id = ?", field.ID).Updates(update_data).Error; err != nil {
			common.ErrorResp(c, err, -1)
			return
		}
		common.SuccessResp(c)
		return
	}

	add_data := &model.AdminMediaInstance{
		Name:       field.Name,
		MediaType:  field.MediaType,
		Mark:       field.Mark,
		HashLife:   field.HashLife,
		Status:     true,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	if err := db.GetDb().Create(add_data).Error; err != nil {
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
