package app

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"mgo/embed"
	"mgo/internal/app/handles"
	handle_admin "mgo/internal/app/handles/admin"
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

	// Build template set with directory-aware names (e.g., "install/index.tmpl")
	// so that we can reference templates across multiple directories explicitly.
	tpl := template.New("").Delims("{[", "]}").Funcs(funcMap)

	for _, name := range embed.TemplatesAllNames("templates") {
		// Trim the leading "templates/" so template names are like "install/index.tmpl"
		short := strings.TrimPrefix(name, "templates/")
		content, err := embed.Templates.ReadFile(name)
		if err != nil {
			panic(err)
		}
		if _, err := tpl.New(short).Parse(string(content)); err != nil {
			panic(err)
		}
	}

	r.SetHTMLTemplate(tpl)
}

// 后台/backstage
func initRuoteAdmin(r *gin.Engine) {
	fmt.Println("conf.Web.AdminPath:", conf.Web.AdminPath)
	backstage := r.Group(conf.Web.AdminPath)
	backstage.GET("/login", handle_admin.LoginPage)
	backstage.POST("/login", handle_admin.PostLogin)
}

func initRuoteInstall(r *gin.Engine) {
	r.GET("/install", install.HomePage)
	r.POST("/install_step1", install.PostInstallStep1)
}

func initRuote(r *gin.Engine) {
	// Static files from embedded filesystem subdir "static"
	staticFS, err := fs.Sub(embed.Static, "static")
	if err != nil {
		panic(err)
	}
	r.StaticFS("/static", http.FS(staticFS))
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	initRuoteAdmin(r)
	initRuoteInstall(r)
	r.GET("/", handles.Home)
}

func Run() {
	r := gin.New()

	// 初始化 session 存储
	store := cookie.NewStore([]byte("mgo"))
	r.Use(sessions.Sessions("app", store))

	// if conf.App.Debug {
	r.Use(gin.Logger())
	// }

	r.Use(gin.Recovery())

	initTemp(r)
	initRuote(r)

	r.Run(fmt.Sprintf(":%d", conf.Web.HTTPPort))
}
