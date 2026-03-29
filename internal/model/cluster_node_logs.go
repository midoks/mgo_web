package model

import (
	"time"
)

type ClusterNodeLogs struct {
	ID          int64     `json:"id" gorm:"primaryKey"`                    // unique key
	Day         int64     `json:"day" gorm:"index"`                        // day
	Description string    `json:"description"`                             // description
	NodeID      int64     `json:"node_id" gorm:"index" binding:"required"` // node_id
	Level       string    `json:"level"`                                   // level
	Tag         string    `json:"tag"`                                     // tag
	IsRead      bool      `json:"is_read" gorm:"index"`                    // is_read
	CreateTime  time.Time `json:"create_time"`                             // create_time
}
