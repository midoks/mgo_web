package app

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"mgo/internal/conf"
)

func Run() {
	r := gin.New()

	r.Run(fmt.Sprintf(":%d", conf.Web.HTTPPort))
}
