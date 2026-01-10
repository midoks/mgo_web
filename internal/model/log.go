package model

import (
	"time"
)

type Log struct {
	ID         int64     `json:"id" gorm:"primaryKey"` // unique key
	Uid        int64     `json:"uid"`                  // uid
	Ip         string    `json:"ip"`                   // ip
	Content    string    `json:"content"`              // content
	CreateTime time.Time `json:"create_time"`          // create_time
	UpdateTime time.Time `json:"update_time"`          // update_time
}
