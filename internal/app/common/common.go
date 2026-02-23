package common

import (
	"context"
	"encoding/json"

	// "fmt"
	"net/http"
	// "strings"

	"mgo/internal/conf"
	"mgo/internal/db"

	// "mgo/internal/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Resp[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

type PageResp struct {
	Content interface{} `json:"content"`
	Total   int64       `json:"total"`
}

type LayuiResp[T any] struct {
	Code  int    `json:"code"`
	Count int64  `json:"count"`
	Msg   string `json:"msg"`
	Data  T      `json:"data"`
}

func CommonVer(c *gin.Context) map[string]interface{} {
	data := map[string]interface{}{
		"title":   "MGOWEB",
		"version": conf.App.Version,
	}

	session := sessions.Default(c)
	username := session.Get("username")
	data["login_name"] = username

	login_user_id := session.Get("user_id")
	var admin_id int64
	switch v := login_user_id.(type) {
	case int64:
		admin_id = v
	case int:
		admin_id = int64(v)
	case uint:
		admin_id = int64(v)
	case uint64:
		admin_id = int64(v)
	}

	data["HasPerm"] = func(code string) bool {
		if admin_id == 0 {
			return false
		}

		ok, err := db.HasAdminPermission(admin_id, code)
		if err != nil {
			return false
		}
		return ok
	}

	data["admin_path"] = conf.Web.AdminPath
	// Build Menus filtered by user's auth codes, super_admin bypass
	menus := GetMenus()
	if admin_id != 0 {
		if u, err := db.GetAdminById(admin_id); err == nil {
			if !u.SuperAdmin {
				allowed := ParseAuthCodes(u.Auth)
				menus = FilterMenusByCodes(menus, allowed)
			}
		}
	}
	data["Menus"] = menus
	data["CurrentPath"] = c.Request.URL.Path
	data["ActiveMenu"] = FindMenuCodeByPath(c.Request.URL.Path, conf.Web.AdminPath)
	return data
}

func ToJson(v interface{}) (d string) {
	rdata, _ := json.MarshalIndent(v, "", "  ")
	return string(rdata)
}

// ErrorResp is used to return error response
// @param l: if true, log error
func ErrorResp(c *gin.Context, err error, code int, l ...bool) {
	ErrorWithDataResp(c, err, code, nil, l...)
}

func hidePrivacy(msg string) string {
	// for _, r := range conf.PrivacyReg {
	// 	msg = r.ReplaceAllStringFunc(msg, func(s string) string {
	// 		return strings.Repeat("*", len(s))
	// 	})
	// }
	return msg
}

func ErrorWithDataResp(c *gin.Context, err error, code int, data interface{}, l ...bool) {
	if len(l) > 0 && l[0] {
		if conf.App.Debug {
			log.Errorf("%+v", err)
		} else {
			log.Errorf("%v", err)
		}
	}
	c.JSON(200, Resp[interface{}]{
		Code: code,
		Msg:  hidePrivacy(err.Error()),
		Data: data,
	})
	c.Abort()
}

func ErrorStrResp(c *gin.Context, str string, code int, l ...bool) {
	if len(l) != 0 && l[0] {
		log.Error(str)
	}
	c.JSON(200, Resp[interface{}]{
		Code: code,
		Msg:  hidePrivacy(str),
		Data: nil,
	})
	c.Abort()
}

func SuccessResp(c *gin.Context, data ...interface{}) {
	if len(data) == 0 {
		c.JSON(200, Resp[interface{}]{
			Code: 200,
			Msg:  "success",
			Data: nil,
		})
		return
	}
	c.JSON(200, Resp[interface{}]{
		Code: 200,
		Msg:  "success",
		Data: data[0],
	})
}

func SuccessLayuiMsgResp(c *gin.Context, msg string, data ...interface{}) {
	if len(data) == 0 {
		c.JSON(200, LayuiResp[interface{}]{
			Code: 0,
			Msg:  msg,
			Data: nil,
		})
		return
	}

	c.JSON(200, LayuiResp[interface{}]{
		Code: 0,
		Msg:  msg,
		Data: data[0],
	})
}

func SuccessLayuiResp(c *gin.Context, count int64, msg string, data ...interface{}) {
	if len(data) == 0 {
		c.JSON(200, LayuiResp[interface{}]{
			Code:  0,
			Count: count,
			Msg:   msg,
			Data:  nil,
		})
		return
	}

	c.JSON(200, LayuiResp[interface{}]{
		Code:  0,
		Count: count,
		Msg:   msg,
		Data:  data[0],
	})
}

func GetHttpReq(ctx context.Context) *http.Request {
	if c, ok := ctx.(*gin.Context); ok {
		return c.Request
	}
	return nil
}
