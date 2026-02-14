package model

import (
	"time"
)

type AdminMediaGroup struct {
	ID         int64     `json:"id" gorm:"primaryKey"`                  // unique key
	Name       string    `json:"name" gorm:"unique" binding:"required"` // name
	Status     bool      `json:"status"`                                // status
	CreateTime time.Time `json:"create_time"`                           // create_time
	UpdateTime time.Time `json:"update_time"`                           // update_time
}
