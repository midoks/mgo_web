package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"mgo/internal/conf"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("logged_in")
		uid := session.Get("user_id")
		// fmt.Println(uid)

		var uidInt int64
		if v, ok := uid.(int64); ok {
			uidInt = v
		}

		if user == nil && uidInt < 1 {
			path := conf.Web.AdminPath
			if !strings.HasPrefix(path, "/") {
				path = "/" + path
			}
			c.Redirect(http.StatusFound, fmt.Sprintf("%s/login", path))
			c.Abort()
			return
		}
		c.Next()
	}
}
