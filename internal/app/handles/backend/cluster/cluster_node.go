package cluster

import (
	"encoding/json"
	"errors"
	"strconv"

	// "fmt"
	"net/http"
	// "strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
	"mgo/internal/db"
	"mgo/internal/model"
	tools "mgo/internal/utils"
)

func parseClusterNodeIpArray(ipJson string) ([]form.ClusterNodeIpAddr, error) {
	if ipJson == "" {
		return nil, nil
	}
	var ipArray []form.ClusterNodeIpAddr
	if err := json.Unmarshal([]byte(ipJson), &ipArray); err != nil {
		return nil, errors.New("invalid ip_addresses_json format: " + ipJson)
	}
	return ipArray, nil
}

func GetNodeSubMenu() []form.SubMenu {
	menu := []form.SubMenu{
		{
			Number: 1,
			Name:   "节点看板",
			Link:   "clusters/node/boards",
		},
		{
			Number: 2,
			Name:   "节点详情",
			Link:   "clusters/node/details",
		},
		{
			Number: 3,
			Name:   "运行日志",
			Link:   "clusters/node/logs",
		},
		{
			Number: 4,
			Name:   "安装节点",
			Link:   "clusters/node/install",
		},
		{
			Number: 5,
			Name:   "节点设置",
			Link:   "clusters/node/settings",
		},
	}
	return menu
}

func Node(c *gin.Context) {
	data := common.CommonVer(c)
	c.HTML(http.StatusOK, "backend/cluster/node.tmpl", data)
}

func CreateNode(c *gin.Context) {
	method := strings.ToUpper(c.Request.Method)
	if method == "POST" {
		PostCreateNode(c)
		return
	}
	data := common.CommonVer(c)
	data["submenu"] = GetSubMenu()
	data["cluster_id"] = c.Query("cluster_id")
	c.HTML(http.StatusOK, "backend/cluster/node/create.tmpl", data)
}

func NodeBoards(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetNodeSubMenu()
	data["node_id"] = c.Query("node_id")
	data["cluster_id"] = c.Query("cluster_id")
	c.HTML(http.StatusOK, "backend/cluster/node/boards.tmpl", data)
}

func NodeDatail(c *gin.Context) {
	data := common.CommonVer(c)
	data["submenu"] = GetNodeSubMenu()
	data["node_id"] = c.Query("node_id")
	data["cluster_id"] = c.Query("cluster_id")

	node_id := c.Query("node_id")
	node_idint, _ := strconv.ParseInt(node_id, 10, 64)
	node_data, _ := db.GetClusterNodeByID(node_idint)
	data["Data"] = node_data

	c.HTML(http.StatusOK, "backend/cluster/node/details.tmpl", data)
}

func NodeList(c *gin.Context) {
	var field form.ClusterNodeList
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}
	result, count, err := db.GetClusterNodeListByClusterID(field.ClusterID, field.Page.Page, field.Page.Limit)
	if err != nil {
		common.ErrorResp(c, err, -2)
		return
	}
	common.SuccessLayuiResp(c, count, "ok", result)
}

func PostCreateNode(c *gin.Context) {
	var field form.ClusterCreateNode
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	if field.Ip == "" {
		common.ErrorResp(c, errors.New("ip address cannot be empty!"), -1)
		return
	}

	secret := tools.RandString(32)
	unique_id := tools.RandString(32)

	nodeip := &model.ClusterNode{
		Name:        field.Name,
		Ip:          field.Ip,
		ClusterID:   field.ClusterID,
		IsInstalled: false,
		Secret:      secret,
		UniqueID:    unique_id,
		CreateTime:  time.Now().Unix(),
		UpdateTime:  time.Now().Unix(),
	}

	if err := db.GetDb().Create(nodeip).Error; err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	if field.IpAddressesJson != "" {
		_, err := parseClusterNodeIpArray(field.IpAddressesJson)
		if err != nil {
			common.ErrorResp(c, err, -1)
			return
		}
	}

	common.SuccessResp(c)
}

func PostDeleteNode(c *gin.Context) {
	var field form.ID
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	err := db.ClusterNodeDeleteByID(nil, field.ID)
	if err == nil {
		common.SuccessResp(c)
		return
	}
	common.ErrorResp(c, err, -1)

}
