package db

import (
	"github.com/pkg/errors"
	// "gorm.io/gorm"

	"mgo/internal/model"
	// utils "mgo/internal/utils"
)

func GetServerList(page, size int) ([]model.Server, int64, error) {
	adminM := db.Model(&model.Server{})
	var count int64
	if err := adminM.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get server count")
	}

	var list []model.Server
	if err := db.Order(columnName("id")).Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return list, count, nil
}

func GetServerById(id int64) (*model.Server, error) {
	var u model.Server
	if err := db.First(&u, id).Error; err != nil {
		return nil, errors.Wrapf(err, "failed get admin")
	}
	return &u, nil
}

func GetServerByIp(ip string) (*model.Server, error) {
	info := model.Server{Ip: ip}
	if err := db.Where(info).First(&info).Error; err != nil {
		return nil, errors.Wrapf(err, "failed find admin")
	}
	return &info, nil
}

func ServerUpdateEmail(id int64, email string) error {
	return db.Model(&model.Server{ID: id}).Update("email", email).Error
}

func ServerUpdateTel(id int64, tel string) error {
	return db.Model(&model.Server{ID: id}).Update("tel", tel).Error
}

func ServerDeleteById(id int64) error {
	var d model.Admin
	return db.Where("id = ?", id).Delete(&d).Error
}
