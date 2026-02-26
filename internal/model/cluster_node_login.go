package model

import (
	"time"
)

type ClusterNodeLogin struct {
	ID         int64     `json:"id" gorm:"primaryKey"` // unique key
	Name       string    `json:"name"`                 // name
	NodeID     int64     `json:"node_id"`              // node_id
	Params     string    `json:"ip" gorm:"unique"`     // params
	Status     bool      `json:"status"`               // status
	CreateTime time.Time `json:"create_time"`          // create_time
	UpdateTime time.Time `json:"update_time"`          // update_time
}
