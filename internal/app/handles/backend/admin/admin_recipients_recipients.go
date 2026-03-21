package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/db"
)

func RecipientsRecipientsDetails(c *gin.Context) {
	id := c.Query("id")
	idint, _ := strconv.ParseInt(id, 10, 64)
	recipient_data, _ := db.GetAdminRecipientsByID(idint)

	data := common.CommonVer(c)
	data["submenu"] = GetRecipientsSubMenu()

	data["id"] = id
	data["Data"] = recipient_data
	c.HTML(http.StatusOK, "backend/admin/recipients/recipients_details.tmpl", data)
}

func RecipientsRecipientsUpdate(c *gin.Context) {
	id := c.Query("id")
	idint, _ := strconv.ParseInt(id, 10, 64)
	recipient_data, _ := db.GetAdminRecipientsByID(idint)

	data := common.CommonVer(c)
	data["id"] = id
	data["Data"] = recipient_data

	cluster_list, _, _ := db.GetClusterList(1, 100)
	data["ClusterList"] = cluster_list

	admin_list, _, _ := db.GetAdminList(1, 100)
	data["AdminList"] = admin_list

	recipients_list, _, _ := db.GetAdminRecipientsInstancesList(1, 100)
	data["RecipientsList"] = recipients_list

	recipients_cluster_related_list, _ := db.GetAdminRecipientsClusterRelatedByRecipientID(idint)
	data["RecipientsClusterRelated"] = recipients_cluster_related_list

	c.HTML(http.StatusOK, "backend/admin/recipients/recipients_update.tmpl", data)
}

func RecipientsRecipientsTest(c *gin.Context) {
	id := c.Query("id")
	idint, _ := strconv.ParseInt(id, 10, 64)
	recipient_data, _ := db.GetAdminRecipientsByID(idint)

	data := common.CommonVer(c)
	data["id"] = id
	data["Data"] = recipient_data
	c.HTML(http.StatusOK, "backend/admin/recipients/recipients_test.tmpl", data)
}
