package db

import (
	"time"

	"mgo/internal/model"
	// "github.com/pkg/errors"
	// "gorm.io/gorm"
)

func ClusterNodeIpaddrSoftDeleteByID(id int64) error {
	if err := db.Model(&model.ClusterNodeIpaddr{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted":  1,
			"update_time": time.Now().Unix(),
		}).Error; err != nil {
		return err
	}
	return nil
}

func ClusterNodeIpaddrDeleteByID(id int64) error {
	var d model.ClusterNodeIpaddr
	return db.Where("id = ?", id).Delete(&d).Error
}
