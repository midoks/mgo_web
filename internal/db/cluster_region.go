package db

import (
	"time"

	"mgo/internal/model"

	"github.com/pkg/errors"
)

func GetClusterRegionList(page, size int) ([]model.ClusterRegion, int64, error) {
	cluster := db.Model(&model.ClusterRegion{})
	var count int64
	if err := cluster.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get cluster region")
	}

	var list []model.ClusterRegion
	if err := db.Order(columnName("id")).Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return list, count, nil
}

func AddClusterRegion(name string, mark string) error {
	data := &model.ClusterRegion{
		Name: name,
		Mark: mark,
	}

	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	if err := errors.WithStack(db.Create(data).Error); err != nil {
		return err
	}
	return nil
}

func UpdateClusterRegion(name string, mark string, id int64) error {
	data := &model.ClusterRegion{
		Name: name,
		Mark: mark,
	}

	data.UpdateTime = time.Now()
	if err := db.Model(&model.ClusterRegion{}).
		Where("id = ?", id).
		Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func GetClusterRegionById(id int64) (*model.ClusterRegion, error) {
	var data model.ClusterRegion
	if err := db.First(&data, id).Error; err != nil {
		return nil, errors.Wrapf(err, "failed get cluster region")
	}
	return &data, nil
}

func ClusterRegionDeleteById(id int64) error {
	var d model.ClusterRegion
	return db.Where("id = ?", id).Delete(&d).Error
}

func ClusterRegionsTriggerStatus(id int64) error {
	var data model.ClusterRegion
	if err := db.First(&data, id).Error; err != nil {
		return errors.Wrapf(err, "failed get cluster region")
	}

	var status int
	if data.Status > 0 {
		status = 0
	} else {
		status = 1
	}

	data.UpdateTime = time.Now()
	data.Status = status

	if err := db.Model(&model.ClusterRegion{}).
		Where("id = ?", id).
		Updates(&data).Error; err != nil {
		return err
	}
	return nil
}
