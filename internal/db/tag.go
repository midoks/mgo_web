package db

import (
	"github.com/pkg/errors"

	"mgo/internal/model"
)

func GetTagList(page, size int) ([]model.Tag, int64, error) {
	tagM := db.Model(&model.Tag{})
	var count int64
	if err := tagM.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get tag count")
	}

	var list []model.Tag
	if err := db.Order(columnName("id")).Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return list, count, nil
}

func GetTagById(id int64) (*model.Tag, error) {
	var u model.Tag
	if err := db.First(&u, id).Error; err != nil {
		return nil, errors.Wrapf(err, "failed get tag")
	}
	return &u, nil
}

func GetTagByName(name string) (*model.Tag, error) {
	info := model.Tag{Name: name}
	if err := db.Where(info).First(&info).Error; err != nil {
		return nil, errors.Wrapf(err, "failed find tag")
	}
	return &info, nil
}

func CreateTag(d *model.Tag) error {
	return db.Create(d).Error
}

func UpdateTag(u *model.Tag) error {
	if err := db.Model(u).Updates(u).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func TagDeleteById(id int64) error {
	var d model.Tag
	if err := db.Where("id = ?", id).Delete(&d).Error; err != nil {
		return errors.Wrapf(err, "failed delete tag")
	}
	return nil
}
