package model

import (
// "fmt"
// "time"
)

type Admin struct {
	ID       int64  `json:"id" gorm:"primaryKey"`                      // unique key
	Username string `json:"username" gorm:"unique" binding:"required"` // username
	Password string `json:"password"`                                  // password
	Salt     string `json:"salt"`                                      // unique salt
}
