package cluster

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
	"mgo/internal/db"
	"mgo/internal/model"
	"mgo/internal/ssh"
)

// 安装状态结构体
type InstallStatus struct {
	NodeID    int64     `json:"node_id"`
	Status    string    `json:"status"`   // pending, running, success, failed
	Progress  int       `json:"progress"` // 0-100
	Message   string    `json:"message"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// 安装状态管理
var (
	installStatusMap   = make(map[int64]*InstallStatus)
	installStatusMutex sync.RWMutex
)

// 获取安装状态
func getInstallStatus(nodeID int64) *InstallStatus {
	installStatusMutex.RLock()
	defer installStatusMutex.RUnlock()
	return installStatusMap[nodeID]
}

// 设置安装状态
func setInstallStatus(nodeID int64, status string, progress int, message string) {
	installStatusMutex.Lock()
	defer installStatusMutex.Unlock()

	if _, exists := installStatusMap[nodeID]; !exists {
		installStatusMap[nodeID] = &InstallStatus{
			NodeID:    nodeID,
			StartTime: time.Now(),
		}
	}

	installStatusMap[nodeID].Status = status
	installStatusMap[nodeID].Progress = progress
	installStatusMap[nodeID].Message = message

	if status == "success" || status == "failed" {
		installStatusMap[nodeID].EndTime = time.Now()
	}
}

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
		fmt.Println(node_data)
		fmt.Println(node_ssh_data)

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
		if err := db.GetDb().Model(&model.ClusterNode{}).Where("id = ?", field.ID).Update("is_installed", field.IsInstalled).Error; err != nil {
			common.ErrorResp(c, err, -1)
			return
		}
		common.SuccessResp(c)
		return
	}
	common.ErrorResp(c, errors.New("node_id error?"), -1)
}

// 开始安装 - 触发安装上传（异步）
func PostNodeInstallDone(c *gin.Context) {
	var field form.ClusterNodeDone
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}
	if field.ID <= 0 {
		common.ErrorResp(c, errors.New("node_id error?"), -1)
		return
	}

	// 检查节点是否存在
	_, err := db.GetClusterNodeByID(field.ID)
	if err != nil {
		common.ErrorResp(c, errors.New("node not found"), -1)
		return
	}

	// 检查是否正在安装
	status := getInstallStatus(field.ID)
	if status != nil && status.Status == "running" {
		common.ErrorResp(c, errors.New("installation is already running"), -1)
		return
	}

	// 设置初始状态
	setInstallStatus(field.ID, "running", 0, "开始安装...")

	// 异步执行安装
	go func(nodeID int64) {
		defer func() {
			if r := recover(); r != nil {
				setInstallStatus(nodeID, "failed", 0, fmt.Sprintf("安装过程发生异常: %v", r))
			}
		}()

		// 执行安装过程
		executeInstallation(nodeID)
	}(field.ID)

	// 返回成功，告知客户端安装已开始
	common.SuccessResp(c)
}

// 执行安装过程
func executeInstallation(nodeID int64) {
	// 获取节点SSH登录信息
	node_login_data, err := db.GetClusterNodeLoginByNodeID(nodeID)
	if err != nil {
		setInstallStatus(nodeID, "failed", 0, "节点SSH登录信息未找到")
		return
	}

	// 获取SSH参数
	login_params, err := node_login_data.GetParams()
	if err != nil {
		setInstallStatus(nodeID, "failed", 0, "获取SSH参数失败")
		return
	}

	// 如果使用了SSH认证，获取SSH认证信息
	var ssh_data *model.ClusterSsh
	if login_params.SshID > 0 {
		ssh_data, err = db.GetClusterSshByID(login_params.SshID)
		if err != nil {
			setInstallStatus(nodeID, "failed", 0, "获取SSH认证信息失败")
			return
		}
	}

	// 准备SSH配置
	ssh_config := ssh.Config{
		Host:    login_params.Host,
		Port:    login_params.Port,
		Timeout: 30 * time.Second,
	}

	if ssh_data != nil {
		ssh_config.User = ssh_data.Username
		ssh_config.Password = ssh_data.Password
		ssh_config.PrivateKeyPEM = []byte(ssh_data.Privatekey)
		ssh_config.PrivateKeyPass = []byte(ssh_data.PrivatekeyPass)
	} else {
		// 如果没有使用SSH认证，使用节点登录信息中的用户名和密码
		ssh_config.User = "root"
		ssh_config.Password = ""
	}

	// 更新状态
	setInstallStatus(nodeID, "running", 10, "正在连接SSH服务器...")

	// 创建SSH客户端
	ssh_client, err := ssh.New(ssh_config)
	if err != nil {
		setInstallStatus(nodeID, "failed", 0, fmt.Sprintf("连接SSH服务器失败: %v", err))
		return
	}
	defer ssh_client.Close()

	// 上传文件
	local_file := filepath.Join("deploy", "mgo_web")
	remote_file := "/tmp/mgo_web"

	// 检查本地文件是否存在
	if _, err := os.Stat(local_file); os.IsNotExist(err) {
		setInstallStatus(nodeID, "failed", 0, "deploy/mgo_web 文件不存在")
		return
	}

	// 更新状态
	setInstallStatus(nodeID, "running", 20, "正在上传文件到服务器...")

	// 上传文件
	err = ssh_client.Upload(local_file, remote_file, 0755, func(written int64, total int64) {
		progress := int(float64(written)/float64(total)*60) + 20 // 20-80%
		setInstallStatus(nodeID, "running", progress, fmt.Sprintf("正在上传文件: %.2f%%", float64(written)/float64(total)*100))
	})
	if err != nil {
		setInstallStatus(nodeID, "failed", 0, fmt.Sprintf("上传文件失败: %v", err))
		return
	}

	// 更新状态
	setInstallStatus(nodeID, "running", 80, "正在执行安装命令...")

	// 在远程服务器上执行安装命令
	install_cmd := fmt.Sprintf("chmod +x %s && %s install", remote_file, remote_file)
	stdout, stderr, err := ssh_client.Run(install_cmd)
	if err != nil {
		setInstallStatus(nodeID, "failed", 0, fmt.Sprintf("执行安装命令失败: %v, stderr: %s", err, stderr))
		return
	}

	// 更新状态
	setInstallStatus(nodeID, "running", 90, "正在配置节点...")

	fmt.Printf("Install output: %s\n", stdout)

	// 安装完成
	setInstallStatus(nodeID, "success", 100, "安装成功完成")
}

// 获取安装状态
func GetNodeInstallStatus(c *gin.Context) {
	nodeIDStr := c.Query("node_id")
	nodeID, err := strconv.ParseInt(nodeIDStr, 10, 64)
	if err != nil {
		common.ErrorResp(c, errors.New("invalid node_id"), -1)
		return
	}

	status := getInstallStatus(nodeID)
	if status == nil {
		// 如果没有安装记录，返回默认状态
		status = &InstallStatus{
			NodeID:   nodeID,
			Status:   "pending",
			Progress: 0,
			Message:  "未开始安装",
		}
	}

	common.SuccessResp(c, status)
}
