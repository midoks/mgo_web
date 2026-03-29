package logs

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
)

// 上报日志
func LogsAdd(c *gin.Context) {
	// 打印请求头信息
	fmt.Println("=== Request Headers ===")
	for key, values := range c.Request.Header {
		for _, value := range values {
			fmt.Printf("header|%s: %s\n", key, value)
		}
	}
	fmt.Println("=======================")

	// 获取原始POST数据
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
	} else {
		fmt.Println("Raw POST data:", string(body))
		c.Request.Body = io.NopCloser(bytes.NewReader(body)) // 重置请求体，以便后续绑定
	}

	var field form.ApiLogs
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	fmt.Println("field:", field)
	common.SuccessResp(c)
}
