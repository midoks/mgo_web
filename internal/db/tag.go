package db

import (
	"github.com/pkg/errors"
	// "gorm.io/gorm"

	"mgo/internal/model"
	// utils "mgo/internal/utils"
)

func GetTagList(page, size int) ([]model.Server, int64, error) {
	serverM := db.Model(&model.Server{})
	var count int64
	if err := serverM.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get server count")
	}

	var list []model.Server
	if err := db.Order(columnName("id")).Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return list, count, nil
}

func GetTagById(id int64) (*model.Server, error) {
	var u model.Server
	if err := db.First(&u, id).Error; err != nil {
		return nil, errors.Wrapf(err, "failed get admin")
	}
	return &u, nil
}

func GetTagByIp(ip string) (*model.Server, error) {
	info := model.Server{Ip: ip}
	if err := db.Where(info).First(&info).Error; err != nil {
		return nil, errors.Wrapf(err, "failed find admin")
	}
	return &info, nil
}

func TagUpdateEmail(id int64, email string) error {
	return db.Model(&model.Server{ID: id}).Update("email", email).Error
}

func TagUpdateTel(id int64, tel string) error {
	return db.Model(&model.Server{ID: id}).Update("tel", tel).Error
}

func TagDeleteById(id int64) error {
	var d model.Admin
	return db.Where("id = ?", id).Delete(&d).Error
}
