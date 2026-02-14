package model

import (
	"time"
)

type AdminMediaInstance struct {
	ID         int64     `json:"id" gorm:"primaryKey"`                  // unique key
	Name       string    `json:"name" gorm:"unique" binding:"required"` // name
	MediaType  string    `json:"media_type"`                            // media_type
	IsOn       string    `json:"is_on"`                                 // is_on
	HashLife   int64     `json:"hash_life"`                             // hash_life
	Params     string    `json:"params"`                                // params
	Rate       string    `json:"rate"`                                  // rate
	Mark       string    `json:"mark"`                                  // mark
	Status     bool      `json:"status"`                                // status
	CreateTime time.Time `json:"create_time"`                           // create_time
	UpdateTime time.Time `json:"update_time"`                           // update_time
}
