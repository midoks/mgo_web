package admin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
	"mgo/internal/db"
	"mgo/internal/model"
)

func GetRecipientsSubMenu() []form.SubMenu {
	menu := []form.SubMenu{
		{
			Number: 1,
			Name:   "接收人",
			Link:   "admin/recipients",
		},
		{
			Number: 2,
			Name:   "接收人分组",
			Link:   "admin/recipients/groups",
		},
		{
			Number: 3,
			Name:   "媒介",
			Link:   "admin/recipients/instances",
		},
		{
			Number: 4,
			Name:   "发送记录",
			Link:   "admin/recipients/logs",
		},
		{
			Number: 5,
			Name:   "任务队列",
			Link:   "admin/recipients/tasks",
		},
	}
	return menu
}

// 通知媒介
func Recipients(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetRecipientsSubMenu()
	c.HTML(http.StatusOK, "backend/admin/recipients/index.tmpl", data)
}

func RecipientsAdd(c *gin.Context) {
	data := common.CommonVer(c)

	cluster_list, _, _ := db.GetClusterList(1, 100)
	data["ClusterList"] = cluster_list

	admin_list, _, _ := db.GetAdminList(1, 100)
	data["AdminList"] = admin_list

	recipients_list, _, _ := db.GetAdminRecipientsInstancesList(1, 100)
	data["RecipientsList"] = recipients_list

	c.HTML(http.StatusOK, "backend/admin/recipients/add.tmpl", data)
}

func PostRecipientsAdd(c *gin.Context) {
	var field form.AdminRecipients
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	common_data := &model.AdminRecipients{
		AdminID:     field.AdminID,
		MediaID:     field.MediaID,
		Status:      field.Status,
		Mark:        field.Mark,
		RecipientID: field.RecipientID,
		ClusterID:   field.ClusterID,
		UpdateTime:  time.Now(),
	}

	if field.ID > 0 {
		if err := db.GetDb().Model(&model.AdminRecipients{}).Where("id = ?", field.ID).Updates(common_data).Error; err != nil {
			common.ErrorResp(c, err, -1)
			return
		}
		common.SuccessResp(c)
		return
	}
	common_data.Status = true
	common_data.CreateTime = time.Now()
	if err := db.GetDb().Create(common_data).Error; err != nil {
		common.ErrorResp(c, err, -1)
		return
	}
	common.SuccessResp(c)
}

func RecipientsList(c *gin.Context) {
	result, count, _ := db.GetAdminRecipientsList(1, 10)
	common.SuccessLayuiResp(c, count, "ok", result)
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
