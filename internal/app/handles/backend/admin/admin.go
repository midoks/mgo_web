package admin

import (
	// "encoding/json"
	"errors"
	// "fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
	"mgo/internal/db"
	"mgo/internal/model"
	utils "mgo/internal/utils"
	// "mgo/internal/op"
)

func Home(c *gin.Context) {
	data := common.CommonVer(c)
	c.HTML(http.StatusOK, "backend/admin/index.tmpl", data)
}

// 通知媒介 -- E

func Add(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)

	admin_data, _ := db.GetAdminById(idInt)
	if admin_data == nil {
		admin_data = &model.Admin{}
	}
	auth := []string{}
	authMap := map[string]bool{}
	if admin_data.Auth != "" {
		auth = strings.Split(admin_data.Auth, ",")
		for _, code := range auth {
			authMap[code] = true
		}
	}

	data := common.CommonVer(c)
	data["Data"] = admin_data
	data["AuthMap"] = authMap
	c.HTML(http.StatusOK, "backend/admin/add.tmpl", data)
}

func PostAdd(c *gin.Context) {
	var f form.AdminAdd
	if err := c.ShouldBind(&f); err != nil {
		common.ErrorResp(c, err, 0)
		return
	}
	f.Auth = c.PostFormMap("auth")

	super_admin := false
	if f.SuperAdmin == "on" {
		super_admin = true
	}

	allow_login := false
	if f.AllowLogin == "on" {
		allow_login = true
	}

	codes := []string{}
	for k, v := range f.Auth {
		if v == "on" {
			codes = append(codes, k)
		}
	}
	codesStr := strings.Join(codes, ",")
	if f.ID > 0 {
		db.UpdateAdmin(f.ID, f.Username, f.Password, f.FullName, codesStr, allow_login, super_admin)
	} else {
		db.AddAdmin(f.Username, f.Password, f.FullName, codesStr, allow_login, super_admin)
	}
	common.SuccessResp(c)
}

func Edit(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)

	admin_data, _ := db.GetAdminById(idInt)

	data := common.CommonVer(c)
	data["Data"] = admin_data
	c.HTML(http.StatusOK, "backend/admin/edit.tmpl", data)
}

func List(c *gin.Context) {
	result, count, _ := db.GetAdminList(1, 10)
	common.SuccessLayuiResp(c, count, "ok", result)
}

func PostEdit(c *gin.Context) {
	var field form.AdminEdit
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, 0)
		return
	}

	add_data := &model.Admin{
		Username: field.Username,
		Password: field.Password,
	}

	if field.ID > 0 {

		if field.Password != "" {
			db.AdminUpdatePass(field.ID, field.Password)
		}

		if field.Tel != "" {
			db.AdminUpdateTel(field.ID, field.Tel)
		}

		if field.Email != "" {
			db.AdminUpdateEmail(field.ID, field.Email)
		}

		common.SuccessResp(c)
		return
	}

	if add_data.Password != "" {
		salt := utils.RandString(16)
		add_data.Salt = salt
		add_data.Password = model.TwoHashPwd(add_data.Password, salt)
	}
	add_data.CreateTime = time.Now()
	add_data.UpdateTime = time.Now()

	err := db.CreateAdmin(add_data)
	if err == nil {
		common.SuccessResp(c)
		return
	}

	common.ErrorResp(c, err, 0)
}

func AdminTriggerStatus(c *gin.Context) {
	var field form.ID
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	err := db.AdminTriggerStatus(field.ID)
	if err == nil {
		common.SuccessResp(c)
		return
	}
	common.ErrorResp(c, err, -1)
}

func Delete(c *gin.Context) {
	var field form.ID
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	if field.ID == 1 {
		common.ErrorResp(c, errors.New("the admin cannot delete!"), -1)
		return
	}

	if err := db.AdminDeleteById(field.ID); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}
	common.SuccessResp(c)
}
