package model

import (
	"time"
)

type Cluster struct {
	ID            int64     `json:"id" gorm:"primaryKey"`                  // unique key
	Name          string    `json:"name" gorm:"unique" binding:"required"` // name
	NodeNum       int64     `json:"node_num"`                              // node num
	NodeNumOnline int64     `json:"node_num_line"`                         // node num online
	CreateTime    time.Time `json:"create_time"`                           // create_time
	UpdateTime    time.Time `json:"update_time"`                           // update_time
}
