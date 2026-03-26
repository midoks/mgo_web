package form

type ClusterNodeDone struct {
	ID int64 `form:"id"`
}

type ClusterNodeUpdateStatus struct {
	ID          int64  `form:"id"`
	IsInstalled string `form:"is_installed"`
}

type ClusterNodeSettings struct {
	ID   int64  `form:"id"`
	Name string `form:"name"`
}

type ClusterNodeLoginAdd struct {
	ID     int64  `form:"id"`
	Host   string `form:"host"`
	Port   int    `form:"port"`
	SshID  int64  `form:"ssh_id"`
	NodeID int64  `form:"node_id"`
}
