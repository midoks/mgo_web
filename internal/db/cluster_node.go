package db

import (
	// "time"

	"github.com/pkg/errors"
	// "gorm.io/gorm"

	"mgo/internal/model"
	// utils "mgo/internal/utils"
)

func GetClusterNodeList(page, size int) ([]model.ClusterNode, int64, error) {
	cluster := db.Model(&model.ClusterNode{})
	var count int64
	if err := cluster.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get cluster group")
	}

	var list []model.ClusterNode
	if err := db.Order(columnName("id")).Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return list, count, nil
}

func GetClusterNodeListByClusterID(cluster_id int64, page, size int) ([]model.ClusterNode, int64, error) {
	cluster := db.Model(&model.ClusterNode{})
	var count int64
	if err := cluster.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get cluster node list")
	}

	var list []model.ClusterNode
	if err := db.Order(columnName("id")).Where("cluster_id =?", cluster_id).Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return list, count, nil
}

func GetClusterNodeByID(id int64) (*model.ClusterNode, error) {
	var data model.ClusterNode
	if err := db.First(&data, id).Error; err != nil {
		return nil, errors.Wrapf(err, "failed get cluster node")
	}
	return &data, nil
}

func ClusterNodeDeleteByID(id int64) error {
	var d model.ClusterNode
	return db.Where("id = ?", id).Delete(&d).Error
}
