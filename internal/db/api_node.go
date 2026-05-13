package db

import (
	"fmt"

	"github.com/pkg/errors"

	"mgo/internal/model"
)

func GetApiNodeList(page, size int) ([]model.ApiNode, int64, error) {
	mm := db.Model(&model.ApiNode{})
	var count int64
	if err := mm.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get api node count")
	}

	var list []model.ApiNode
	if err := db.Order(columnName("id")).Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return list, count, nil
}

func GetApiNodeByID(id int64) (*model.ApiNode, error) {
	var u model.ApiNode
	if err := db.First(&u, id).Error; err != nil {
		return nil, errors.Wrapf(err, "failed get api log")
	}
	return &u, nil
}

func GetApiNodeAddr() string {
	var node model.ApiNode

	if err := db.Where("status = ? AND is_primary = ?", true, true).First(&node).Error; err == nil {
		return fmt.Sprintf("%s:%d", node.Domain, node.Port)
	}

	if err := db.Where("status = ?", true).Order("`order` asc").First(&node).Error; err == nil {
		return fmt.Sprintf("%s:%d", node.Domain, node.Port)
	}

	return "http://127.0.0.1:9292"
}
