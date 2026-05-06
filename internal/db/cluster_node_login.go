package db

import (
	"mgo/internal/model"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
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

func ClusterNodeLoginDeleteByID(tx *gorm.DB, id int64) error {
	if tx == nil {
		tx = db
	}
	var d model.ClusterNodeLogin
	return tx.Where("id = ?", id).Delete(&d).Error
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
		Where("status = ? AND node_id IN (?)", 1, sub).
		Select("CAST(JSON_EXTRACT(params, '$.ssh_id') AS UNSIGNED) as ssh_id, COUNT(*) AS c").
		Group("ssh_id").
		Having("ssh_id > 0").
		Order("c DESC").
		Limit(3)

	err := qb.Find(&rows).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ids := make([]int64, 0, len(rows))
	for _, r := range rows {
		if r.SshID > 0 {
			ids = append(ids, r.SshID)
		}
	}
	return ids, nil
}

func ClusterNodeLoginAddOrUpdate(nodeID int64, host string, port int, sshID int64) error {
	common_data := &model.ClusterNodeLogin{
		Name:       "ssh",
		NodeID:     nodeID,
		UpdateTime: time.Now().Unix(),
	}

	common_data.SetParams(model.ClusterNodeLoginParams{
		Host:  host,
		Port:  port,
		SshID: sshID,
	})

	var existing model.ClusterNodeLogin
	err := db.Where("node_id = ?", nodeID).First(&existing).Error
	if err == nil && existing.ID > 0 {
		return db.Model(&model.ClusterNodeLogin{}).Where("id = ?", existing.ID).Updates(common_data).Error
	}

	common_data.Status = true
	common_data.CreateTime = time.Now().Unix()
	return db.Create(common_data).Error
}
