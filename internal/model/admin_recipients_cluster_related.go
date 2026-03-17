package model

import (
	"time"
)

type AdminRecipientsClusterRelated struct {
	ID          int64     `json:"id" gorm:"primaryKey"` // unique key
	RecipientID string    `json:"recipient_id"`         // recipient_id
	ClusterID   int64     `json:"cluster_id"`           // cluster_id
	Status      bool      `json:"status"`               // status
	CreateTime  time.Time `json:"create_time"`          // create_time
	UpdateTime  time.Time `json:"update_time"`          // update_time
}
