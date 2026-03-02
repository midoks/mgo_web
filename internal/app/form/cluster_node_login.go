package form

type ClusterNodeLoginAdd struct {
	ID     int64  `form:"id"`
	Host   string `form:"host"`
	Port   int    `form:"port"`
	SshID  int64  `form:"ssh_id"`
	NodeID int64  `form:"node_id"`
}
