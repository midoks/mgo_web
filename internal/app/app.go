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
		"HasPrefix": func(s, prefix string) bool {
			return strings.HasPrefix(s, prefix)
		},
		//是子菜单或当前菜单
		"IsSubOrEq": func(base, menu string) bool {
			if base == menu {
				return true
			}
			endp := strings.Replace(base, menu, "", 1)
			endp = strings.TrimPrefix(endp, "/")
			return !strings.Contains(endp, "/")
		},
		"Contains": func(s, substr string) bool {
			return strings.Contains(s, substr)
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
	backstage_admin.GET("/admin/list", backend_admin.List)
	backstage_admin.GET("/admin/details", backend_admin.Details)
	backstage_admin.GET("/admin/update", backend_admin.Update)
	backstage_admin.POST("/admin/delete", backend_admin.Delete)
	backstage_admin.POST("/admin/trigger_status", backend_admin.AdminTriggerStatus)

	// 管理员 - 通知
	backstage_admin.GET("/admin/recipients", backend_admin.Recipients)
	backstage_admin.GET("/admin/recipients/list", backend_admin.RecipientsList)
	backstage_admin.POST("/admin/recipients/delete", backend_admin.RecipientsDelete)
	backstage_admin.GET("/admin/recipients/add", backend_admin.RecipientsAdd)
	backstage_admin.POST("/admin/recipients/add", backend_admin.PostRecipientsAdd)
	backstage_admin.GET("/admin/recipients/groups", backend_admin.RecipientsGroups)
	backstage_admin.GET("/admin/recipients/groups/list", backend_admin.RecipientsGroupsList)
	backstage_admin.GET("/admin/recipients/groups/select", backend_admin.RecipientsGroupsSelect)
	backstage_admin.GET("/admin/recipients/groups/add", backend_admin.RecipientsGroupsAdd)
	backstage_admin.POST("/admin/recipients/groups/add", backend_admin.PostRecipientsGroupsAdd)
	backstage_admin.POST("/admin/recipients/groups/delete", backend_admin.PostRecipientsGroupsDelete)
	backstage_admin.GET("/admin/recipients/instances", backend_admin.RecipientsInstances)
	backstage_admin.GET("/admin/recipients/instances/list", backend_admin.RecipientsInstancesList)
	backstage_admin.GET("/admin/recipients/instances/add", backend_admin.RecipientsInstancesAdd)
	backstage_admin.POST("/admin/recipients/instances/add", backend_admin.PostRecipientsInstancesAdd)
	backstage_admin.GET("/admin/recipients/instances/details", backend_admin.RecipientsInstancesDetails)
	backstage_admin.GET("/admin/recipients/instances/update", backend_admin.RecipientsInstancesUpdate)
	backstage_admin.GET("/admin/recipients/instances/test", backend_admin.RecipientsInstancesTest)
	backstage_admin.POST("/admin/recipients/instances/test", backend_admin.PostRecipientsInstancesTest)
	backstage_admin.POST("/admin/recipients/instances/delete", backend_admin.RecipientsInstancesDelete)

	backstage_admin.GET("/admin/recipients/tasks", backend_admin.RecipientsTasks)
	backstage_admin.GET("/admin/recipients/logs", backend_admin.RecipientsLogs)

	// 边缘节点
	backstage_admin.GET("/clusters", backend_cluster.Home)
	backstage_admin.POST("/clusters/create", backend_cluster.PostCreate)
	backstage_admin.GET("/clusters/list", backend_cluster.List)
	backstage_admin.POST("/clusters/delete", backend_cluster.Delete)

	backstage_admin.GET("/clusters/cluster/boards", backend_cluster.ClusterBoards)
	backstage_admin.GET("/clusters/cluster/list", backend_cluster.ClusterList)

	// 边缘节点 - 节点
	backstage_admin.GET("/clusters/node", backend_cluster.Node)
	backstage_admin.GET("/clusters/ipaddr", backend_cluster.Node)
	backstage_admin.GET("/clusters/create", backend_cluster.Create)
	backstage_admin.GET("/clusters/select/ip", backend_cluster.SelectIp)
	backstage_admin.GET("/clusters/select/region", backend_cluster.SelectRegion)
	backstage_admin.GET("/clusters/select/groups", backend_cluster.SelectGroups)
	backstage_admin.GET("/clusters/cluster/create_node", backend_cluster.CreateNode)
	backstage_admin.POST("/clusters/cluster/create_node", backend_cluster.PostCreateNode)
	backstage_admin.POST("/clusters/cluster/delete_node", backend_cluster.PostDeleteNode)
	backstage_admin.GET("/clusters/cluster/node_list", backend_cluster.NodeList)

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

	//边缘节点 - 认证
	backstage_admin.GET("/clusters/ssh", backend_cluster.ClusterSsh)
	backstage_admin.GET("/clusters/ssh/list", backend_cluster.ClusterSshList)
	backstage_admin.GET("/clusters/ssh/details", backend_cluster.ClusterSshDetails)
	backstage_admin.GET("/clusters/ssh/test", backend_cluster.ClusterSshTest)
	backstage_admin.GET("/clusters/ssh/add", backend_cluster.ClusterSshAdd)
	backstage_admin.GET("/clusters/ssh/create", backend_cluster.ClusterSshCreate)
	backstage_admin.POST("/clusters/ssh/create", backend_cluster.PostClusterSshCreate)
	backstage_admin.GET("/clusters/ssh/update", backend_cluster.ClusterSshUpdate)

	// 日志审计
	backstage_admin.GET("/log", backend_log.Home)
	backstage_admin.GET("/log/list", backend_log.List)
	backstage_admin.GET("/log/settings", backend_log.Settings)
	backstage_admin.POST("/log/settings", backend_log.PostSettting)
	backstage_admin.GET("/log/clean", backend_log.Clean)
	backstage_admin.POST("/log/clean", backend_log.PostLogClean)
	backstage_admin.POST("/log/delete", backend_log.Delete)

	// 系统设置
	backstage_admin.GET("/system/settings", backend_system.Home)
	backstage_admin.GET("/system/settings/profile", backend_system.Profile)
	backstage_admin.POST("/system/settings/profile", backend_system.PostProfile)

	backstage_admin.GET("/system/settings/login", backend_system.Login)
	backstage_admin.POST("/system/settings/login", backend_system.PostLogin)
	backstage_admin.GET("/system/settings/login/logs", backend_system.LoginLogs)
	backstage_admin.GET("/system/settings/login/logs/list", backend_system.LoginLogsList)
	backstage_admin.GET("/system/database", backend_system.Database)
	backstage_admin.GET("/system/db", backend_system.Db)

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
