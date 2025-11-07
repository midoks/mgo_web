package model

import (
	"fmt"
	"time"

	utils "mgo/internal/utils"
)

const StaticHashSalt = "mgo"

type Admin struct {
	ID         int64     `json:"id" gorm:"primaryKey"`                      // unique key
	Username   string    `json:"username" gorm:"unique" binding:"required"` // username
	Password   string    `json:"password"`                                  // password
	Salt       string    `json:"salt"`                                      // salt
	Status     int       `json:"status"`                                    // status
	Tel        string    `json:"tel"`                                       // tel
	Email      string    `json:"email"`                                     // email
	CreateTime time.Time `json:"create_time"`                               // create_time
	UpdateTime time.Time `json:"update_time"`                               // update_time
}

func StaticHash(password string) string {
	return utils.HashData(utils.SHA256, []byte(fmt.Sprintf("%s-%s", password, StaticHashSalt)))
}

func HashPwd(static string, salt string) string {
	return utils.HashData(utils.SHA256, []byte(fmt.Sprintf("%s-%s", static, salt)))
}

func TwoHashPwd(password string, salt string) string {
	return HashPwd(StaticHash(password), salt)
}
