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

func GetAdminList(page, size int) ([]model.Admin, int64, error) {
	adminM := db.Model(&model.Admin{})
	var count int64
	if err := adminM.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get server count")
	}

	var list []model.Admin
	if err := db.Order(columnName("id")).Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return list, count, nil
}

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
