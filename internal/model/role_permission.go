package model

type RolePermission struct {
	ID           int64 `json:"id" gorm:"primaryKey"`
	RoleID       int64 `json:"role_id" gorm:"index"`
	PermissionID int64 `json:"permission_id" gorm:"index"`
}
