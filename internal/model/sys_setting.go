package model

import (
	"encoding/json"
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
	AllowedManualDelete bool   `json:"allowed_manual_delete"`  // allowed manual delete
	AllowedManual       bool   `json:"allowed_manual"`         // allowed manual
	SaveDay             int64  `json:"save_day"`               // save day
	MaxCapacityLimit    string `json:"max_capacity_limit"`     // max capacity limit
	AllowModClearConfig string `json:"allow_mod_clear_config"` // allow mod clear config
}

func (a *SysSetting) SetEmailParams(p SysSettingLogValue) error {
	b, err := json.Marshal(p)
	if err != nil {
		return err
	}
	a.Value = string(b)
	return nil
}

func (a *SysSetting) GetEmailParams() (SysSettingLogValue, error) {
	var p SysSettingLogValue
	if a.Value == "" {
		return p, nil
	}
	err := json.Unmarshal([]byte(a.Value), &p)
	return p, err
}
