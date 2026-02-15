package model

import (
	"time"
)

type AdminRecipients struct {
	ID          int64     `json:"id" gorm:"primaryKey"` // unique key
	AdminID     int64     `json:"admin_id"`             // admin_id
	MediaID     int64     `json:"media_id"`             // media_id
	RecipientID string    `json:"recipient_id"`         // recipient_id
	ClusterID   int64     `json:"cluster_id"`           // cluster_id
	Mark        string    `json:"mark"`                 // mark
	Status      bool      `json:"status"`               // status
	CreateTime  time.Time `json:"create_time"`          // create_time
	UpdateTime  time.Time `json:"update_time"`          // update_time
}
