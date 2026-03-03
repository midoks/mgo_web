package model

import (
	"time"
)

type ClusterNode struct {
	ID               int64     `json:"id" gorm:"primaryKey"`                             // unique key
	ClusterID        int64     `json:"cluster_id" gorm:"index"`                          // cluster_id
	Name             string    `json:"name"`                                             // name
	Ip               string    `json:"ip" gorm:"unique;index" binding:"required"`        // ip
	AllowPublic      bool      `json:"allow_public" binding:"required"`                  // allow_public
	AllowHealthCheck bool      `json:"allow_health_check" binding:"required"`            // allow_health_check
	IsInstalled      bool      `json:"is_installed"`                                     // is_installed
	NodeInfo         string    `json:"node_info"`                                        // node_info
	Secret           string    `json:"secret" gorm:"unique;index" binding:"required"`    // secret
	UniqueID         string    `json:"unique_id" gorm:"unique;index" binding:"required"` // unique_id
	Mark             string    `json:"mark"`                                             // mark
	Status           bool      `json:"status" gorm:"index"`                              // status
	CreateTime       time.Time `json:"create_time"`                                      // create_time
	UpdateTime       time.Time `json:"update_time"`                                      // update_time
}
