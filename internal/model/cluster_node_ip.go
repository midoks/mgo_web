package model

import (
	"time"
)

type ClusterNodeIp struct {
	ID               int64     `json:"id" gorm:"primaryKey"`                  // unique key
	ClusterID        int64     `json:"cluster_id"`                            // cluster_id
	Name             string    `json:"name"`                                  // name
	Ip               string    `json:"ip" gorm:"unique" binding:"required"`   // ip
	AllowPublic      bool      `json:"allow_public" binding:"required"`       // allow_public
	AllowHealthCheck bool      `json:"allow_health_check" binding:"required"` // allow_health_check
	Mark             string    `json:"mark"`                                  // mark
	Status           bool      `json:"status"`                                // status
	CreateTime       time.Time `json:"create_time"`                           // create_time
	UpdateTime       time.Time `json:"update_time"`                           // update_time
}
