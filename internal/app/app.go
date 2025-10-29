package app

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/handles"
	"mgo/internal/conf"
)

func Run() {
	r := gin.New()
	r.Use(gin.Recovery())

	if conf.App.Debug {
		r.Use(gin.Logger())
	}

	r.Any("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.GET("/", handles.HomePage)

	r.Run(fmt.Sprintf(":%d", conf.Web.HTTPPort))
}
