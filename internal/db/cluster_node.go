package db

import (
	"time"

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

func AddClusterNodeGroup(name string, clusterId int64) error {
	data := &model.ClusterGroup{
		Name:      name,
		ClusterId: clusterId,
	}

	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	if err := errors.WithStack(db.Create(data).Error); err != nil {
		return err
	}
	return nil
}

func UpdateClusterNodeGroup(name string, id int64) error {
	data := &model.ClusterGroup{
		Name: name,
	}

	data.UpdateTime = time.Now()
	if err := db.Model(&model.ClusterGroup{}).
		Where("id = ?", id).
		Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func GetClusterNodeGroupById(id int64) (*model.ClusterGroup, error) {
	var data model.ClusterGroup
	if err := db.First(&data, id).Error; err != nil {
		return nil, errors.Wrapf(err, "failed get cluster group")
	}
	return &data, nil
}

func ClusterNodeDeleteById(id int64) error {
	var d model.ClusterNode
	return db.Where("id = ?", id).Delete(&d).Error
}
