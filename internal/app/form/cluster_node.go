package form

type ClusterNodeQuery struct {
	ID int64 `form:"id"`
	Page
}

type ClusterNodeDone struct {
	ID int64 `form:"id"`
}

type ClusterNodeUpdateStatus struct {
	ID          int64  `form:"id"`
	IsInstalled string `form:"is_installed"`
}

type ClusterNodeIpAddr struct {
	Ip             int64  `form:"ip"`
	AllowPublic    string `form:"allow_public"`
	CanHealthCheck string `form:"can_health_check"`
	IsOn           string `form:"is_on"`
	Description    string `form:"description"`
}

type ClusterNodeSettings struct {
	ID              int64               `form:"id"`
	Name            string              `form:"name"`
	IpAddressesJson []ClusterNodeIpAddr `form:"ip_addresses_json"`
}

type ClusterNodeLoginAdd struct {
	ID     int64  `form:"id"`
	Host   string `form:"host"`
	Port   int    `form:"port"`
	SshID  int64  `form:"ssh_id"`
	NodeID int64  `form:"node_id"`
}
