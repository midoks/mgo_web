package form

type ClusterCreate struct {
	Name string `form:"name"`
}

type ClusterSubMenu struct {
	Number int64  `form:"number"`
	Name   string `form:"name"`
	Link   string `form:"link"`
}

type ClusterGroupAdd struct {
	ID        string `form:"id"`
	Name      string `form:"name"`
	ClusterId int64  `form:"cluster_id"`
}

type CreateNode struct {
	Name string `form:"name"`
	Ip   string `form:"ip"`
}
