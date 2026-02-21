package db

import (
	"mgo/internal/model"

	"github.com/pkg/errors"
)

func GetClusterSshList(page, size int) ([]model.ClusterSsh, int64, error) {
	cluster := db.Model(&model.ClusterSsh{})
	var count int64
	if err := cluster.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get cluster region")
	}

	var list []model.ClusterSsh
	if err := db.Order(columnName("id")).Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return list, count, nil
}

func GetClusterSshById(id int64) (*model.ClusterSsh, error) {
	var data model.ClusterSsh
	if err := db.First(&data, id).Error; err != nil {
		return nil, errors.Wrapf(err, "failed get cluster region")
	}
	return &data, nil
}

func ClusterSshDeleteById(id int64) error {
	var d model.ClusterSsh
	return db.Where("id = ?", id).Delete(&d).Error
}
