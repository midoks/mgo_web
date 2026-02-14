package db

import (
	"github.com/pkg/errors"

	"mgo/internal/model"
)

func GetAdminRecipientsList(page, size int) ([]model.AdminMediaInstance, int64, error) {
	adminM := db.Model(&model.AdminMediaInstance{})
	var count int64
	if err := adminM.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get recipients data list")
	}

	var list []model.AdminMediaInstance
	if err := db.Order(columnName("id")).Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return list, count, nil
}

func GetAdminRecipientById(id int64) (*model.AdminMediaInstance, error) {
	var u model.AdminMediaInstance
	if err := db.First(&u, id).Error; err != nil {
		return nil, errors.Wrapf(err, "failed get recipients data")
	}
	return &u, nil
}

func AdminRecipientsDeleteById(id int64) error {
	var d model.AdminMediaInstance
	return db.Where("id = ?", id).Delete(&d).Error
}
