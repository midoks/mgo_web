package db

import (
	"github.com/pkg/errors"
	"mgo/internal/model"
)

func GetSysSettingByCode(code string) (*model.SysSetting, error) {
	var u model.SysSetting
	if err := db.Where("code = ?", code).First(&u).Error; err != nil {
		return nil, errors.Wrapf(err, "failed get sys setting")
	}
	return &u, nil
}

func SysSettingDeleteByCode(code string) error {
	var d model.SysSetting
	return db.Where("code = ?", code).Delete(&d).Error
}
