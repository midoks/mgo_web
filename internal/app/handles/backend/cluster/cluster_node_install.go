package cluster

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	// "strings"
	// "time"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
	"mgo/internal/db"
	"mgo/internal/model"
)

func NodeInstall(c *gin.Context) {
	node_id := c.Query("node_id")
	node_idint, _ := strconv.ParseInt(node_id, 10, 64)

	data := common.CommonVer(c)
	data["submenu"] = GetNodeSubMenu()
	data["node_id"] = node_id
	data["cluster_id"] = c.Query("cluster_id")

	if node_id != "" {
		node_data, _ := db.GetClusterNodeByID(node_idint)
		data["Data"] = node_data

		node_ssh_data, _ := db.GetClusterNodeLoginByNodeID(node_idint)
		data["SshData"] = node_ssh_data
	}

	c.HTML(http.StatusOK, "backend/cluster/node/install.tmpl", data)
}

func PostNodeInstallUpdateStatus(c *gin.Context) {
	var field form.ClusterNodeUpdateStatus
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}
	if field.ID > 0 {
		if err := db.GetDb().Model(&model.ClusterNode{ID: field.ID}).Update("is_installed", field.IsInstalled).Error; err != nil {
			common.ErrorResp(c, err, -1)
			return
		}
		common.SuccessResp(c)
	}
	common.ErrorResp(c, errors.New("node_id error?"), -1)
}
