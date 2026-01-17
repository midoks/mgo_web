package db

import (
	"time"

	"mgo/internal/model"

	"github.com/pkg/errors"
)

func GetLogList(page, size int) ([]model.Log, int64, error) {
	serverM := db.Model(&model.Log{})
	var count int64
	if err := serverM.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get server count")
	}

	var list []model.Log
	if err := db.Order(columnName("id")).Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return list, count, nil
}

func GetLogById(id int64) (*model.Log, error) {
	var u model.Log
	if err := db.First(&u, id).Error; err != nil {
		return nil, errors.Wrapf(err, "failed get admin")
	}
	return &u, nil
}

func LogDeleteById(id int64) error {
	var d model.Admin
	return db.Where("id = ?", id).Delete(&d).Error
}

func AddLog(uid int64, content string) error {
	var u model.Log
	u.Uid = uid
	u.Content = content
	u.CreateTime = time.Now()

	return errors.WithStack(db.Create(&u).Error)
}
