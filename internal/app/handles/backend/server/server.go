package server

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/db"
	"mgo/internal/model"
	utils "mgo/internal/utils"
	// "mgo/internal/op"
)

func Home(c *gin.Context) {
	data := common.CommonVer()
	data["PageIsServer"] = true

	c.HTML(http.StatusOK, "backend/server/index.tmpl", data)
}

func Edit(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)

	admin_data, _ := db.GetAdminById(idInt)

	data := common.CommonVer()
	data["Data"] = admin_data
	c.HTML(http.StatusOK, "backend/server/edit.tmpl", data)
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
