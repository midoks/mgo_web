package form

type ClusterCreate struct {
	Name string `form:"name"`
}

type ClusterSubMenu struct {
	Number int64  `form:"number"`
	Name   string `form:"name"`
	Link   string `form:"link"`
}
