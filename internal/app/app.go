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
	backend_system "mgo/internal/app/handles/backend/system"
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

	backstage_admin.GET("/admin/add", backend_admin.Add)
	backstage_admin.POST("/admin/add", backend_admin.PostAdd)
	backstage_admin.GET("/admin/edit", backend_admin.Edit)
	backstage_admin.POST("/admin/edit", backend_admin.PostEdit)
	backstage_admin.GET("/admin/list", backend_admin.List)
	backstage_admin.POST("/admin/delete", backend_admin.Delete)
	backstage_admin.POST("/admin/trigger_status", backend_admin.AdminTriggerStatus)

	// 管理员 - 通知
	backstage_admin.GET("/admin/recipients", backend_admin.Recipients)
	backstage_admin.GET("/admin/recipients/groups", backend_admin.RecipientsGroups)
	backstage_admin.GET("/admin/recipients/instances", backend_admin.RecipientsInstances)
	backstage_admin.GET("/admin/recipients/instances/add", backend_admin.RecipientsInstancesAdd)

	// 边缘节点
	backstage_admin.GET("/clusters", backend_cluster.Home)
	backstage_admin.GET("/clusters/list", backend_cluster.List)
	backstage_admin.POST("/clusters/delete", backend_cluster.Delete)

	backstage_admin.GET("/clusters/cluster/boards", backend_cluster.ClusterBoards)
	backstage_admin.GET("/clusters/cluster/list", backend_cluster.ClusterList)

	//边缘节点 - 认证
	backstage_admin.GET("/clusters/grants", backend_cluster.ClusterGrants)

	// 边缘节点 - 节点
	backstage_admin.GET("/clusters/node", backend_cluster.Node)
	backstage_admin.GET("/clusters/create", backend_cluster.Create)
	backstage_admin.GET("/clusters/select/region", backend_cluster.SelectRegion)
	backstage_admin.GET("/clusters/select/groups", backend_cluster.SelectGroups)
	backstage_admin.GET("/clusters/cluster/create_node", backend_cluster.ClusterCreateNode)

	// 边缘节点 - 分组
	backstage_admin.GET("/clusters/cluster/groups", backend_cluster.ClusterGroups)
	backstage_admin.GET("/clusters/cluster/groups_add", backend_cluster.ClusterGroupsAdd)
	backstage_admin.GET("/clusters/cluster/groups_list", backend_cluster.ClusterGroupsList)
	backstage_admin.POST("/clusters/cluster/groups_add", backend_cluster.PostClusterGroupsAdd)
	backstage_admin.POST("/clusters/cluster/groups_delete", backend_cluster.ClusterGroupsDelete)

	backstage_admin.GET("/clusters/cluster/install", backend_cluster.ClusterInstall)
	backstage_admin.GET("/clusters/cluster/settings", backend_cluster.ClusterSettings)
	backstage_admin.GET("/clusters/cluster/delete", backend_cluster.ClusterDelete)

	// 边缘节点 - 区域设置
	backstage_admin.GET("/clusters/regions", backend_cluster.ClusterRegions)
	backstage_admin.GET("/clusters/regions/list", backend_cluster.ClusterRegionsList)
	backstage_admin.GET("/clusters/regions/add", backend_cluster.ClusterRegionsAdd)
	backstage_admin.GET("/clusters/regions/nodes", backend_cluster.ClusterRegionsNodes)
	backstage_admin.POST("/clusters/regions/add", backend_cluster.PostClusterRegionsNodesAdd)
	backstage_admin.POST("/clusters/regions/delete", backend_cluster.ClusterRegionsDelete)
	backstage_admin.POST("/clusters/regions/trigger_status", backend_cluster.ClusterRegionsTriggerStatus)

	backstage_admin.POST("/clusters/create", backend_cluster.PostCreate)

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

	// 系统设置
	backstage_admin.GET("/system/base", backend_system.Home)
	backstage_admin.GET("/system/database", backend_system.Home)

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
