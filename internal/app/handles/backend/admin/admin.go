package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
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

func Recipients(c *gin.Context) {
	data := common.CommonVer(c)
	c.HTML(http.StatusOK, "backend/admin/recipients.tmpl", data)
}

func Add(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)

	admin_data, _ := db.GetAdminById(idInt)

	data := common.CommonVer(c)
	data["Data"] = admin_data
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

	b, _ := json.Marshal(f)
	fmt.Println("field:", string(b))
	fmt.Println("auth:", f.Auth)
	if f.ID > 0 {
		db.UpdateAdmin(f.ID, f.Username, f.Password, f.FullName, allow_login, super_admin)
	} else {
		db.AddAdmin(f.Username, f.Password, f.FullName, allow_login, super_admin)
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
	var f struct {
		Id       int64  `form:"id"`
		Username string `form:"username"`
		Tel      string `form:"Tel"`
		Email    string `form:"email"`
		Password string `form:"password"`
	}

	if err := c.ShouldBind(&f); err != nil {
		common.ErrorResp(c, err, 0)
		return
	}

	d := &model.Admin{
		Username: f.Username,
		Password: f.Password,
	}

	if f.Id > 0 {

		if f.Password != "" {
			db.AdminUpdatePass(f.Id, f.Password)
		}

		if f.Tel != "" {
			db.AdminUpdateTel(f.Id, f.Tel)
		}

		if f.Email != "" {
			db.AdminUpdateEmail(f.Id, f.Email)
		}

		common.SuccessResp(c)
		return
	}

	if d.Password != "" {
		salt := utils.RandString(16)
		d.Salt = salt
		d.Password = model.TwoHashPwd(d.Password, salt)
	}
	d.CreateTime = time.Now()
	d.UpdateTime = time.Now()

	err := db.CreateAdmin(d)
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
	var f struct {
		Id int64 `form:"id"`
	}

	if err := c.ShouldBind(&f); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	if f.Id == 1 {
		common.ErrorResp(c, errors.New("the admin cannot delete!"), -1)
		return
	}

	err := db.AdminDeleteById(f.Id)
	if err == nil {
		common.SuccessResp(c)
		return
	}
	common.ErrorResp(c, err, -1)
}
