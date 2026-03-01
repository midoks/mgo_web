package db

import (
	// "fmt"
	"mgo/internal/model"

	"github.com/pkg/errors"
	// "gorm.io/gorm"
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

func GetClusterNodeLoginByID(id int64) (*model.ClusterNodeLogin, error) {
	var data model.ClusterNodeLogin
	if err := db.Where("id=?", id).First(&data).Error; err != nil {
		return nil, errors.Wrapf(err, "failed get cluster node login")
	}
	return &data, nil
}

func GetClusterNodeLoginByNodeID(node_id int64) (*model.ClusterNodeLogin, error) {
	var data model.ClusterNodeLogin
	if err := db.Where("node_id=?", node_id).First(&data).Error; err != nil {
		return nil, errors.Wrapf(err, "failed get cluster node login")
	}
	return &data, nil
}

func ClusterNodeLoginDeleteById(id int64) error {
	var d model.ClusterNodeLogin
	return db.Where("id = ?", id).Delete(&d).Error
}

func ClusterNodeLoginFindFrequentSshIDs(clusterID int64) ([]int64, error) {
	type row struct {
		SshID int64 `gorm:"column:ssh_id"`
		C     int64 `gorm:"column:c"`
	}
	var rows []row
	sub := db.Model(&model.ClusterNode{}).
		Select("id").
		Where("cluster_id = ?", clusterID)
	qb := db.Model(&model.ClusterNodeLogin{}).
		Where("status = ?", 1).
		Where("node_id IN (?)", sub).
		Select("CAST(JSON_EXTRACT(params, '$.ssh_id') AS INTEGER) as ssh_id, COUNT(*) AS c").
		Group("ssh_id").
		Having("ssh_id > 0").
		Order("c DESC").
		Limit(3)
	// sql := qb.Session(&gorm.Session{DryRun: true}).ToSQL(func(tx *gorm.DB) *gorm.DB {
	// 	return tx.Find(&rows)
	// })
	// fmt.Println("SQL:", sql)
	err := qb.Find(&rows).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var ids []int64
	for _, r := range rows {
		if r.SshID > 0 {
			ids = append(ids, r.SshID)
		}
	}
	return ids, nil
}
