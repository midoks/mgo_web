package db

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"mgo/internal/model"
)

func GetAdminRecipientsList(page, size int) ([]model.AdminRecipients, int64, error) {
	adminM := db.Model(&model.AdminRecipients{})
	var count int64
	if err := adminM.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get recipients data list")
	}

	var list []model.AdminRecipients
	if err := db.Order(columnName("id")).Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return list, count, nil
}

func GetAdminRecipientsByID(id int64) (*model.AdminRecipients, error) {
	var u model.AdminRecipients
	if err := db.First(&u, id).Error; err != nil {
		return nil, errors.Wrapf(err, "failed get recipients data")
	}
	return &u, nil
}

func AdminRecipientsDeleteById(tx *gorm.DB, id int64) error {
	if tx == nil {
		tx = db
	}
	var d model.AdminRecipients
	return tx.Where("id = ?", id).Delete(&d).Error
}
