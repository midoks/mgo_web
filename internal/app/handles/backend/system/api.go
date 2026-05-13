package system

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
	"mgo/internal/db"
	"mgo/internal/model"
)

func Api(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetSysAdvancedSubMenu()
	c.HTML(http.StatusOK, "backend/system/api/index.tmpl", data)
}

func ApiAdd(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetSysAdvancedSubMenu()
	c.HTML(http.StatusOK, "backend/system/api/add.tmpl", data)
}

func ApiDetails(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetSysAdvancedSubMenu()

	id := c.Query("id")
	data["id"] = id

	idInt, _ := strconv.ParseInt(id, 10, 64)
	dbnode_data, err := db.GetApiNodeByID(idInt)
	if err == nil {
		data["Data"] = dbnode_data
	}

	c.HTML(http.StatusOK, "backend/system/api/details.tmpl", data)
}

func ApiClean(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetSysAdvancedSubMenu()

	id := c.Query("id")
	data["id"] = id

	idInt, _ := strconv.ParseInt(id, 10, 64)
	dbnode_data, err := db.GetDbNodeByID(idInt)
	if err == nil {
		data["Data"] = dbnode_data
	}

	c.HTML(http.StatusOK, "backend/system/api/clean.tmpl", data)
}

func ApiLogs(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetSysAdvancedSubMenu()

	id := c.Query("id")
	data["id"] = id

	idInt, _ := strconv.ParseInt(id, 10, 64)
	dbnode_data, err := db.GetDbNodeByID(idInt)
	if err == nil {
		data["Data"] = dbnode_data
	}

	c.HTML(http.StatusOK, "backend/system/api/logs.tmpl", data)
}

func ApiUpdate(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetSysAdvancedSubMenu()

	id := c.Query("id")
	data["id"] = id

	idInt, _ := strconv.ParseInt(id, 10, 64)
	dbnode_data, err := db.GetDbNodeByID(idInt)
	if err == nil {
		data["Data"] = dbnode_data
	}

	c.HTML(http.StatusOK, "backend/system/api/update.tmpl", data)
}

func ApiList(c *gin.Context) {
	var field form.Page
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	result, count, err := db.GetApiNodeList(field.Page, field.Limit)
	if err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	common.SuccessLayuiResp(c, count, "ok", result)
}

func PostApiAdd(c *gin.Context) {
	var field form.ApiNode
	var err error
	if err = c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	common_data := &model.ApiNode{
		Name:       field.Name,
		Domain:     field.Domain,
		Type:       field.Type,
		Order:      0,
		Weigth:     0,
		IsPrimary:  true,
		Status:     field.Status,
		UpdateTime: time.Now().Unix(),
	}

	if field.ID > 0 {
		if err := db.GetDb().Model(&model.ApiNode{}).Where("id = ?", field.ID).Updates(common_data).Error; err != nil {
			common.ErrorResp(c, err, -1)
			return
		}
	} else {
		common_data.CreateTime = time.Now().Unix()
		if err := db.GetDb().Create(common_data).Error; err != nil {
			common.ErrorResp(c, err, -1)
			return
		}
	}
	common.SuccessResp(c)
}
