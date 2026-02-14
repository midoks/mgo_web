package db

import (
	"time"

	"github.com/pkg/errors"
	// "gorm.io/gorm"

	"mgo/internal/model"
	// utils "mgo/internal/utils"
)

func GetClusterNodeIpList(page, size int) ([]model.ClusterNodeIp, int64, error) {
	cluster := db.Model(&model.ClusterNodeIp{})
	var count int64
	if err := cluster.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get cluster group")
	}

	var list []model.ClusterNodeIp
	if err := db.Order(columnName("id")).Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return list, count, nil
}

func GetClusterNodeIpListByClusterID(cluster_id int64, page, size int) ([]model.ClusterNodeIp, int64, error) {
	cluster := db.Model(&model.ClusterNodeIp{})
	var count int64
	if err := cluster.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get cluster group")
	}

	var list []model.ClusterNodeIp
	if err := db.Order(columnName("id")).Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return list, count, nil
}

func AddClusterNodeIpGroup(name string, clusterId int64) error {
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

func UpdateClusterNodeIpGroup(name string, id int64) error {
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

func GetClusterNodeIpGroupById(id int64) (*model.ClusterGroup, error) {
	var data model.ClusterGroup
	if err := db.First(&data, id).Error; err != nil {
		return nil, errors.Wrapf(err, "failed get cluster group")
	}
	return &data, nil
}

func ClusterNodeIpDeleteById(id int64) error {
	var d model.ClusterNodeIp
	return db.Where("id = ?", id).Delete(&d).Error
}
