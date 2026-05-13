package model

type ApiNode struct {
	ID         int64  `json:"id" gorm:"primaryKey"` // unique key
	Name       string `json:"name"`                 // name
	Type       string `json:"type"`                 // type
	Domain     string `json:"domain"`               // domain
	IsPrimary  bool   `json:"is_primary"`           // is_primary
	Port       int64  `json:"port"`                 // port
	Order      int64  `json:"order"`                // order
	Weigth     int64  `json:"weigth"`               // weigth
	Status     bool   `json:"status"`               // status
	Mark       string `json:"mark"`                 // mark
	CreateTime int64  `json:"create_time"`          // create_time
	UpdateTime int64  `json:"update_time"`          // update_time
}
