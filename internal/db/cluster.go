package db

import (
	"github.com/pkg/errors"
	// "gorm.io/gorm"

	"mgo/internal/model"
	// utils "mgo/internal/utils"
)

func GetClusterList(page, size int) ([]model.Cluster, int64, error) {
	cluster := db.Model(&model.Cluster{})
	var count int64
	if err := cluster.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get server count")
	}

	var list []model.Cluster
	if err := db.Order(columnName("id")).Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return list, count, nil
}

func GetClusterById(id int64) (*model.Server, error) {
	var u model.Server
	if err := db.First(&u, id).Error; err != nil {
		return nil, errors.Wrapf(err, "failed get admin")
	}
	return &u, nil
}

func GetClusterByIp(ip string) (*model.Server, error) {
	info := model.Server{Ip: ip}
	if err := db.Where(info).First(&info).Error; err != nil {
		return nil, errors.Wrapf(err, "failed find admin")
	}
	return &info, nil
}

func ClusterDeleteById(id int64) error {
	var d model.Admin
	return db.Where("id = ?", id).Delete(&d).Error
}
