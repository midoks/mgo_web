package install

import (
	// "fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/conf"
	// "mgo/internal/db"
)

func HomePage(c *gin.Context) {
	data := common.CommonVer()
	step := c.Query("step")
	if step == "2" {
		c.HTML(http.StatusOK, "install/step2.tmpl", data)
		return
	}
	c.HTML(http.StatusOK, "install/index.tmpl", data)
}

func Step2Page(c *gin.Context) {
	data := common.CommonVer()
	c.HTML(http.StatusOK, "install/step2.tmpl", data)
}

func PostInstallStep1(c *gin.Context) {
	install_data := make(map[string]string, 0)
	install_data["type"] = c.PostForm("type")
	install_data["hostname"] = c.PostForm("hostname")
	install_data["hostport"] = c.PostForm("hostport")
	install_data["dbname"] = c.PostForm("dbname")
	install_data["username"] = c.PostForm("username")
	install_data["password"] = c.PostForm("password")
	install_data["table_prefix"] = c.PostForm("table_prefix")
	install_data["dbpath"] = c.PostForm("dbpath")

	err := conf.InstallConf(install_data)
	if err != nil {
		common.ErrorStrResp(c, err.Error(), -1)
		return
	}

	// init_account := c.PostForm("account")
	// init_pass := c.PostForm("pass")
	// if conf.Security.InstallLock {
	// 	db.InitDb()
	// 	userdata.InitAdmin(init_account, init_pass)
	// }

	common.SuccessResp(c, gin.H{"token": "安装成功!"})
}
