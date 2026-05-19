package logs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
	// "mgo/internal/conf"
	"mgo/internal/db"
	"mgo/internal/model"
)

func DebugInfo(c *gin.Context) {
	// 获取原始POST数据
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
	} else {
		fmt.Println("Raw POST data:", string(body))
		c.Request.Body = io.NopCloser(bytes.NewReader(body)) // 重置请求体，以便后续绑定
	}
}

// 上报日志
func LogsAdd(c *gin.Context) {

	api_header := c.Request.Header
	unique_id := api_header.Get("X-Node-Id")
	secret := api_header.Get("X-Secret")

	// 尝试获取SQLite数据库文件路径
	sqlDB, _ := db.GetDb().DB()
	if sqlDB != nil {
		fmt.Println("DB Connection OK")
	}

	node_data, err := db.GetClusterNodeByUniqueIdAndSecret(unique_id, secret)
	if err != nil {
		fmt.Println("GetClusterNodeByUniqueIdAndSecret error:", err)
		common.ErrorResp(c, err, -1)
		return
	}

	// DebugInfo(c)

	var field form.ApiLogs
	if err := c.ShouldBind(&field); err != nil {
		fmt.Println("err:", err)
		common.ErrorResp(c, err, -1)
		return
	}

	now := time.Now().Unix()
	fmt.Println("field:", now, field)

	// 解析 Data 字段
	if field.Data != "" {
		var dataMap map[string]interface{}
		if err := json.Unmarshal([]byte(field.Data), &dataMap); err == nil {
			itemType, _ := dataMap["item"].(string)
			valueStr, _ := dataMap["value"].(string)

			if itemType == "sysinfo" && valueStr != "" {
				// 解析 sysinfo 数据
				var nodeInfo model.ClusterNodeNodeInfoParam
				if err := json.Unmarshal([]byte(valueStr), &nodeInfo); err == nil {
					// 更新节点信息
					node_data.SetNodeInfoParams(nodeInfo)
					if err := db.GetDb().Save(node_data).Error; err != nil {
						fmt.Println("Update cluster node error:", err)
					}
				}
			} else {
				fmt.Println("item:", node_data)
			}
		}
	}

	common.SuccessResp(c)
}
