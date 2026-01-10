package model

type Role struct {
	ID   int64  `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"uniqueIndex;size:64"`
	Desc string `json:"desc" gorm:"size:255"`
}
