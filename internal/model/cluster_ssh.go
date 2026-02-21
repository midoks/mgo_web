package model

import (
	"time"
)

type ClusterSsh struct {
	ID             int64     `json:"id" gorm:"primaryKey"`                  // unique key
	Name           string    `json:"name" gorm:"unique" binding:"required"` // name
	Method         int       `json:"method"`                                // method
	Username       string    `json:"username"`                              // username
	Password       string    `json:"password"`                              // password
	PrivateKey     string    `json:"private_key"`                           // private_key
	PrivateKeyPass string    `json:"private_key_pass"`                      // private_key_pass
	Mark           string    `json:"mark"`                                  // mark
	CreateTime     time.Time `json:"create_time"`                           // create_time
	UpdateTime     time.Time `json:"update_time"`                           // update_time
}
