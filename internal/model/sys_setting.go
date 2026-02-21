package model

import (
	"encoding/json"
	"fmt"
	"time"
)

type SysSetting struct {
	ID         int64     `json:"id" gorm:"primaryKey"` // unique key
	Code       string    `json:"code" gorm:"unique"`   // code
	Uid        int64     `json:"uid"`                  // uid
	Value      string    `json:"value"`                // value
	UpdateTime time.Time `json:"update_time"`          // update_time
	CreateTime time.Time `json:"create_time"`          // create_time
}

type SysSettingLogValue struct {
	AllowedManualDelete   bool   `json:"allowed_manual_delete"`    // allowed manual delete
	AllowedManual         bool   `json:"allowed_manual"`           // allowed manual
	SaveDay               int64  `json:"save_day"`                 // save day
	MaxCapacityLimit      int64  `json:"max_capacity_limit"`       // max capacity limit
	MaxCapacityUnit       string `json:"max_capacity_unit"`        // max capacity unit
	AllowedModClearConfig bool   `json:"allowed_mod_clear_config"` // allowed mod clear config
}

func (a *SysSetting) SetLogValue(p SysSettingLogValue) error {
	b, err := json.Marshal(p)
	if err != nil {
		return err
	}
	a.Value = string(b)
	return nil
}

func (a *SysSetting) GetLogValue() (SysSettingLogValue, error) {
	var p SysSettingLogValue
	if a.Value == "" {
		return p, nil
	}

	err := json.Unmarshal([]byte(a.Value), &p)

	fmt.Println("p:", p)
	return p, err
}
