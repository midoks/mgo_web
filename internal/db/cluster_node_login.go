package db

import (
	// "time"

	"github.com/pkg/errors"
	// "gorm.io/gorm"

	"mgo/internal/model"
	// utils "mgo/internal/utils"
)

func GetClusterNodeLoginList(page, size int) ([]model.ClusterNodeLogin, int64, error) {
	cluster := db.Model(&model.ClusterNodeLogin{})
	var count int64
	if err := cluster.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get cluster login")
	}

	var list []model.ClusterNodeLogin
	if err := db.Order(columnName("id")).Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return list, count, nil
}

func ClusterNodeLoginDeleteById(id int64) error {
	var d model.ClusterNodeLogin
	return db.Where("id = ?", id).Delete(&d).Error
}
