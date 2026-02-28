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
		return nil, 0, errors.Wrapf(err, "failed get cluster node login")
	}

	var list []model.ClusterNodeLogin
	if err := db.Order(columnName("id")).Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return list, count, nil
}

func GetClusterNodeByID(id int64) (*model.ClusterNodeLogin, error) {
	var data model.ClusterNodeLogin
	if err := db.Where("id=?", id).First(&data).Error; err != nil {
		return nil, errors.Wrapf(err, "failed get cluster group")
	}
	return &data, nil
}

func GetClusterNodeByNodeID(node_id int64) (*model.ClusterNodeLogin, error) {
	var data model.ClusterNodeLogin
	if err := db.Where("node_id=?", node_id).First(&data).Error; err != nil {
		return nil, errors.Wrapf(err, "failed get cluster group")
	}
	return &data, nil
}

func ClusterNodeLoginDeleteById(id int64) error {
	var d model.ClusterNodeLogin
	return db.Where("id = ?", id).Delete(&d).Error
}

func ClusterNodeLoginFindFrequentSsh(cluster_id int64) {

}
