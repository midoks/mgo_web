package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/db"
	// "mgo/internal/model"
)

func RecipientsRecipientsDetails(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetRecipientsSubMenu()

	id := c.Query("id")
	idint, _ := strconv.ParseInt(id, 10, 64)
	recipient_data, _ := db.GetAdminRecipientsById(idint)

	fmt.Println(recipient_data)
	data["id"] = id
	data["Data"] = recipient_data
	c.HTML(http.StatusOK, "backend/admin/recipients/recipients_details.tmpl", data)
}

func RecipientsRecipientsUpdate(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)
	recipient_data, _ := db.GetAdminRecipientsInstancesById(idInt)

	data := common.CommonVer(c)
	data["id"] = id
	data["Data"] = recipient_data
	c.HTML(http.StatusOK, "backend/admin/recipients/recipients_update.tmpl", data)
}

func RecipientsRecipientsTest(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)
	recipient_data, _ := db.GetAdminRecipientsInstancesById(idInt)

	data := common.CommonVer(c)
	data["id"] = id
	data["Data"] = recipient_data

	c.HTML(http.StatusOK, "backend/admin/recipients/recipients_test.tmpl", data)
}
