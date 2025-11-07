package db

import (
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

func CreateAdmin(u *model.Admin) error {
	return errors.WithStack(db.Create(u).Error)
}
