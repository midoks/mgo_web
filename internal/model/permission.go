package model

type Permission struct {
	ID   int64  `json:"id" gorm:"primaryKey"`
	Code string `json:"code" gorm:"uniqueIndex;size:128"`
	Desc string `json:"desc" gorm:"size:255"`
}
