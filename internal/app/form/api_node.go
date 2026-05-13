package form

type ApiNode struct {
	ID     int64  `form:"id"`
	Name   string `form:"name"`
	Type   string `form:"type"`
	Domain string `form:"domain"`
	Status bool   `form:"status"`
}
