package model

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"mgo/internal/errs"
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

func (u *Admin) ValidatePwdStaticHash(password string) error {
	if password == "" {
		return errors.WithStack(errs.EmptyPassword)
	}

	fmt.Println("login", u.Password, HashPwd(password, u.Salt), password, u.Salt)
	if u.Password != HashPwd(password, u.Salt) {
		return errors.WithStack(errs.WrongPassword)
	}
	return nil
}

func HashPwd(password string, salt string) string {
	return utils.HashData(utils.SHA256, []byte(fmt.Sprintf("%s-%s", password, salt)))
}

func TwoHashPwd(password string, salt string) string {
	fmt.Println("reg", password, HashPwd(password, salt), password, salt)
	return HashPwd(password, salt)
}
