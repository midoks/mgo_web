package model

import (
	"time"
)

type SysSetting struct {
	ID         int64     `json:"id" gorm:"primaryKey"` // unique key
	Code       string    `json:"code"`                 // code
	Uid        int64     `json:"uid"`                  // uid
	Value      string    `json:"value"`                // value
	UpdateTime time.Time `json:"update_time"`          // update_time
	CreateTime time.Time `json:"create_time"`          // create_time
}
