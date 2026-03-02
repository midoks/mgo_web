package cluster

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/entity"
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
	c.HTML(http.StatusOK, "backend/cluster/ssh/index.tmpl", data)
}

func ClusterSshAdd(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetSshSubMenu()

	ssh_id := c.Query("ssh_id")
	data["ssh_id"] = ssh_id
	ssh_idint, _ := strconv.ParseInt(ssh_id, 10, 64)
	ssh_data, _ := db.GetClusterSshByID(ssh_idint)
	data["Data"] = ssh_data

	c.HTML(http.StatusOK, "backend/cluster/ssh/add.tmpl", data)
}

func ClusterSshCreate(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetSshSubMenu()
	c.HTML(http.StatusOK, "backend/cluster/ssh/create.tmpl", data)
}

func ClusterSshDetails(c *gin.Context) {
	ssh_id := c.Query("ssh_id")

	data := common.CommonVer(c)
	data["ssh_id"] = ssh_id

	ssh_idint, _ := strconv.ParseInt(ssh_id, 10, 64)
	ssh_data, _ := db.GetClusterSshByID(ssh_idint)
	data["Data"] = ssh_data

	c.HTML(http.StatusOK, "backend/cluster/ssh/details.tmpl", data)
}

func ClusterSshUpdate(c *gin.Context) {
	ssh_id := c.Query("ssh_id")

	data := common.CommonVer(c)
	data["ssh_id"] = c.Query("ssh_id")

	ssh_idint, _ := strconv.ParseInt(ssh_id, 10, 64)
	ssh_data, _ := db.GetClusterSshByID(ssh_idint)
	data["Data"] = ssh_data
	c.HTML(http.StatusOK, "backend/cluster/ssh/update.tmpl", data)
}

func ClusterSshTest(c *gin.Context) {
	ssh_id := c.Query("ssh_id")

	data := common.CommonVer(c)
	data["ssh_id"] = c.Query("ssh_id")

	ssh_idint, _ := strconv.ParseInt(ssh_id, 10, 64)
	ssh_data, _ := db.GetClusterSshByID(ssh_idint)
	data["Data"] = ssh_data
	c.HTML(http.StatusOK, "backend/cluster/ssh/test.tmpl", data)
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

func ClusterSshSelectList(c *gin.Context) {
	common_data := &entity.ClusterSshEntityList{}
	limit_data, _ := db.GetClusterSshListByLimit(100)
	common_data.List = limit_data

	sugguest_data, _ := db.GetClusterSshListBySuggest(1)
	common_data.Sugguest = sugguest_data
	common.SuccessResp(c, common_data)
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
