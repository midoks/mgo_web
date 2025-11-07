package db

import (
	"fmt"
	// "encoding/base64"
	// "sync"
	// "time"

	"github.com/pkg/errors"

	"mgo/internal/model"
	// "mgo/internal/utils"
)

func GetAdminById(id int64) (*model.Admin, error) {
	var u model.Admin
	if err := db.First(&u, id).Error; err != nil {
		return nil, errors.Wrapf(err, "failed get admin")
	}
	return &u, nil
}

func GetAdminByName(username string) (*model.Admin, error) {
	info := model.Admin{Username: username}
	if err := db.Where(info).First(&info).Error; err != nil {
		return nil, errors.Wrapf(err, "failed find admin")
	}
	return &info, nil
}

func CreateAdmin(u *model.Admin) error {

	fmt.Println(u)
	return errors.WithStack(db.Create(u).Error)
}
