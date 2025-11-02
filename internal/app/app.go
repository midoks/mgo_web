package app

import (
	"fmt"
	"html/template"

	"github.com/gin-gonic/gin"

	"mgo/embed"
	"mgo/internal/app/handles"
	"mgo/internal/app/handles/install"
	"mgo/internal/conf"
)

func initTemp(r *gin.Engine) {
	// Define template functions
	funcMap := template.FuncMap{
		"safe": func(str string) template.HTML {
			return template.HTML(str)
		},
		// Cache-busting token exposed as a function for templates
		"BuildCommit": func() string {
			return conf.BuildCommit
		},
	}

	// Use custom delimiters to avoid conflicts with client-side "{{ }}" templates
	// and apply function map before parsing templates
	tpl := template.Must(
		template.New("").Delims("{[", "]}").Funcs(funcMap).ParseFS(
			embed.Templates,
			"templates/*.tmpl",
			"templates/**/*.tmpl",
		),
	)

	r.SetHTMLTemplate(tpl)
}

func initRuote(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.GET("/install", install.Home)
	r.GET("/", handles.Home)
}

func Run() {
	r := gin.New()

	if conf.App.Debug {
		r.Use(gin.Logger())
	}

	r.Use(gin.Recovery())

	initTemp(r)
	initRuote(r)

	r.Run(fmt.Sprintf(":%d", conf.Web.HTTPPort))
}
