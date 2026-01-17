package model

import (
	"time"
)

type ClusterGroup struct {
	ID         int64     `json:"id" gorm:"primaryKey"`                  // unique key
	Name       string    `json:"name" gorm:"unique" binding:"required"` // name
	CreateTime time.Time `json:"create_time"`                           // create_time
	UpdateTime time.Time `json:"update_time"`                           // update_time
}
