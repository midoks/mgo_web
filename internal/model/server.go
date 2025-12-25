package model

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"mgo/internal/errs"
	utils "mgo/internal/utils"
)

type Server struct {
	ID         int64     `json:"id" gorm:"primaryKey"`                // unique key
	Ip         string    `json:"ip" gorm:"unique" binding:"required"` // ip
	Tags       string    `json:"tags"`                                // tags
	Status     int       `json:"status"`                              // status
	CreateTime time.Time `json:"create_time"`                         // create_time
	UpdateTime time.Time `json:"update_time"`                         // update_time
}
