package db

import (
	"fmt"

	"mgo/internal/app/entity"
	"mgo/internal/model"

	"github.com/pkg/errors"
)

func GetClusterSshList(page, size int) ([]model.ClusterSsh, int64, error) {
	cluster := db.Model(&model.ClusterSsh{})
	var count int64
	if err := cluster.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get cluster ssh")
	}

	var list []model.ClusterSsh
	if err := db.Order(columnName("id")).Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return list, count, nil
}

func GetClusterSshListByLimit(limit int) ([]entity.ClusterSsh, error) {
	var models []model.ClusterSsh
	if err := db.Order(columnName("id")).Limit(limit).Find(&models).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	out := make([]entity.ClusterSsh, 0, len(models))
	for _, m := range models {
		out = append(out, entity.ClusterSsh{
			ID:       m.ID,
			Name:     m.Name,
			Method:   m.Method,
			Username: m.Username,
		})
	}
	return out, nil
}

func GetClusterSshListBySuggest(clusterID int64) ([]entity.ClusterSsh, error) {
	out := []entity.ClusterSsh{}

	ids, err := ClusterNodeLoginFindFrequentSshIDs(clusterID)
	fmt.Println("ids:", ids, err)
	if err != nil {
		return out, err
	}

	if len(ids) == 0 {
		return out, nil
	}

	var models []model.ClusterSsh
	if err := db.Order(columnName("id")).Where("id IN (?)", ids).Limit(3).Find(&models).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	out = make([]entity.ClusterSsh, 0, len(models))
	for _, m := range models {
		out = append(out, entity.ClusterSsh{
			ID:       m.ID,
			Name:     m.Name,
			Method:   m.Method,
			Username: m.Username,
		})
	}
	return out, nil
}

func GetClusterSshById(id int64) (*model.ClusterSsh, error) {
	var data model.ClusterSsh
	if err := db.First(&data, id).Error; err != nil {
		return nil, errors.Wrapf(err, "failed get cluster ssh")
	}
	return &data, nil
}

func ClusterSshDeleteById(id int64) error {
	var d model.ClusterSsh
	return db.Where("id = ?", id).Delete(&d).Error
}
