package db

import (
	"time"

	"github.com/pkg/errors"

	"mgo/internal/model"
)

func GetAdminRecipientsList(page, size int) ([]model.AdminMediaInstance, int64, error) {
	adminM := db.Model(&model.AdminMediaInstance{})
	var count int64
	if err := adminM.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get server count")
	}

	var list []model.AdminMediaInstance
	if err := db.Order(columnName("id")).Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return list, count, nil
}

func AddAdminRecipients(name string) error {
	data := &model.AdminMediaInstance{
		Name: name,
	}

	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	if err := errors.WithStack(db.Create(data).Error); err != nil {
		return err
	}
	return nil
}

func AdminRecipientsDeleteById(id int64) error {
	var d model.AdminMediaInstance
	return db.Where("id = ?", id).Delete(&d).Error
}
