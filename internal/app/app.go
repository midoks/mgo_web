package app

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	// "time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"mgo/embed"
	"mgo/internal/app/handles"
	backend "mgo/internal/app/handles/backend"
	backend_admin "mgo/internal/app/handles/backend/admin"
	backend_cluster "mgo/internal/app/handles/backend/cluster"
	backend_log "mgo/internal/app/handles/backend/log"
	backend_server "mgo/internal/app/handles/backend/server"
	"mgo/internal/app/handles/install"
	"mgo/internal/app/middleware"
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
	// fmt.Println("conf.Web.AdminPath:", conf.Web.AdminPath)
	backstage := r.Group(conf.Web.AdminPath)
	backstage.GET("/login", backend.LoginPage)
	backstage.POST("/login", backend.PostLogin)
	backstage.GET("/logout", backend.LoginOut)

	backstage_admin := backstage.Group("")
	backstage_admin.Use(middleware.AuthRequired())

	// 管理员
	backstage_admin.GET("", backend.HomePage)
	backstage_admin.GET("/index", handles.Home)
	backstage_admin.GET("/admin/index", backend_admin.Home)
	backstage_admin.GET("/admin/recipients", backend_admin.Recipients)
	backstage_admin.GET("/admin/edit", middleware.PermissionRequired("admin.edit"), backend_admin.Edit)
	backstage_admin.POST("/admin/edit", middleware.PermissionRequired("admin.edit"), backend_admin.PostEdit)
	backstage_admin.GET("/admin/list", middleware.PermissionRequired("admin.view"), backend_admin.List)
	backstage_admin.POST("/admin/delete", middleware.PermissionRequired("admin.edit"), backend_admin.Delete)

	// 边缘节点
	backstage_admin.GET("/clusters", backend_cluster.Home)
	backstage_admin.GET("/clusters/list", backend_cluster.List)
	backstage_admin.GET("/clusters/create", backend_cluster.Create)
	backstage_admin.GET("/clusters/node", backend_cluster.Node)

	backstage_admin.POST("/clusters/create", backend_cluster.PostCreate)

	backstage_admin.GET("/clusters/cluster/boards", backend_cluster.ClusterBoards)
	backstage_admin.GET("/clusters/cluster/list", backend_cluster.ClusterList)

	//服务器
	backstage_admin.GET("/server/index", backend_server.Home)
	backstage_admin.GET("/server/list", middleware.PermissionRequired("server.view"), backend_server.List)
	backstage_admin.GET("/server/edit", middleware.PermissionRequired("server.edit"), backend_server.Edit)

	// 日志审计
	backstage_admin.GET("/log", backend_log.Home)
	backstage_admin.GET("/log/list", backend_log.List)

	backstage_admin.GET("/tag/index", backend_server.Home)
	backstage_admin.GET("/tag/list", middleware.PermissionRequired("server.view"), backend_server.List)
	backstage_admin.GET("/tag/edit", middleware.PermissionRequired("server.edit"), backend_server.Edit)

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
	// 	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
	// 		p := param.Path
	// 		if strings.Contains(p, ".js") || strings.Contains(p, ".css") {
	// 			return ""
	// 		}
	// 		if strings.Contains(p, ".woff2") {
	// 			return ""
	// 		}
	// 		return fmt.Sprintf("%s - [%s] \"%s %s %s\" %d %s \"%s\"\n",
	// 			param.ClientIP,
	// 			param.TimeStamp.Format(time.RFC1123),
	// 			param.Method,
	// 			p,
	// 			param.Request.Proto,
	// 			param.StatusCode,
	// 			param.Latency,
	// 			param.ErrorMessage,
	// 		)
	// 	}))
	// }

	// r.Use(gin.Recovery())
	r.SetTrustedProxies(nil)

	initTemp(r)
	initRuote(r)

	r.Run(fmt.Sprintf(":%d", conf.Web.HTTPPort))
}
