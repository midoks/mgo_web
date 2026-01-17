package model

import (
	"time"
)

type Cluster struct {
	ID         int64     `json:"id" gorm:"primaryKey"`                  // unique key
	Name       string    `json:"name" gorm:"unique" binding:"required"` // name
	Num        int64     `json:"num"`                                   // node num
	NumOnline  int64     `json:"num_online"`                            // node num online
	CreateTime time.Time `json:"create_time"`                           // create_time
	UpdateTime time.Time `json:"update_time"`                           // update_time
}
