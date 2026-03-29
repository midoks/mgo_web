package form

type ApiLogs struct {
	NodeId    string `json:"node_id"`   // node_id
	Secret    string `json:"secret"`    // secret
	Version   string `json:"version"`   // 请求版本信息
	Timestamp string `json:"timestamp"` // 时间戳

	// type:
	// sys		-> 系统信息
	// node 	-> 运行信息
	// request 	-> 应用请求
	Type string `json:"type"`
}
