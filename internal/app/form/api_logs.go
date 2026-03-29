package form

type ApiLogs struct {
	NodeId  string `form:"node_id"` // node_id
	Secret  string `form:"secret"`  // secret
	Version string `form:"version"` // 请求版本信息

	// type:
	// sys		-> 系统信息
	// node 	-> 运行信息
	// request 	-> 应用请求
	Type string `form:"type"`
}
