package form

type ClusterNodeLoginAdd struct {
	ID     int64  `form:"id"`
	Host   string `form:"host"`
	Port   string `form:"port"`
	SshID  string `form:"ssh_id"`
	NodeID int64  `form:"node_id"`
}
