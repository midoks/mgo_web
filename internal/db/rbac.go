package db

import (
	"strings"
	"time"

	"github.com/pkg/errors"

	"mgo/internal/model"
)

func InitRBACDefault() {
	var count int64
	db.Model(&model.Role{}).Where("name = ?", "super_admin").Count(&count)
	if count == 0 {
		r := &model.Role{Name: "super_admin", Desc: "超级管理员"}
		_ = db.Create(r).Error
	}
	ensurePermission("admin.view", "查看管理员")
	ensurePermission("admin.edit", "编辑管理员")
	ensurePermission("server.view", "查看边缘节点")
	ensurePermission("server.edit", "编辑边缘节点")

	var super model.Role
	if err := db.Where("name = ?", "super_admin").First(&super).Error; err == nil {
		var perms []model.Permission
		db.Find(&perms)
		for _, p := range perms {
			var cnt int64
			db.Model(&model.RolePermission{}).
				Where("role_id = ? AND permission_id = ?", super.ID, p.ID).
				Count(&cnt)
			if cnt == 0 {
				db.Create(&model.RolePermission{RoleID: super.ID, PermissionID: p.ID})
			}
		}
	}

	var admin model.Admin
	if err := db.First(&admin, 1).Error; err == nil {
		var cnt int64
		db.Model(&model.AdminRole{}).Where("admin_id = ? AND role_id = ?", admin.ID, super.ID).Count(&cnt)
		if cnt == 0 {
			db.Create(&model.AdminRole{AdminID: admin.ID, RoleID: super.ID})
		}
	}
}

func ensurePermission(code, desc string) {
	code = strings.TrimSpace(code)
	if code == "" {
		return
	}
	var count int64
	db.Model(&model.Permission{}).Where("code = ?", code).Count(&count)
	if count == 0 {
		db.Create(&model.Permission{Code: code, Desc: desc})
	}
}

func GetAdminPermissions(adminID int64) ([]string, error) {
	type Row struct {
		Code string
	}
	var rows []Row
	err := db.Table("permissions p").
		Select("p.code").
		Joins("JOIN role_permissions rp ON rp.permission_id = p.id").
		Joins("JOIN admin_roles ar ON ar.role_id = rp.role_id").
		Where("ar.admin_id = ?", adminID).
		Scan(&rows).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	perms := make([]string, 0, len(rows))
	for _, r := range rows {
		perms = append(perms, r.Code)
	}
	return perms, nil
}

func HasAdminPermission(adminID int64, code string) (bool, error) {
	var cnt int64
	err := db.Table("admin_roles ar").
		Joins("JOIN roles r ON r.id = ar.role_id").
		Where("ar.admin_id = ? AND r.name = ?", adminID, "super_admin").
		Count(&cnt).Error
	if err != nil {
		return false, errors.WithStack(err)
	}
	if cnt > 0 {
		return true, nil
	}
	var c int64
	err = db.Table("role_permissions rp").
		Joins("JOIN permissions p ON p.id = rp.permission_id").
		Joins("JOIN admin_roles ar ON ar.role_id = rp.role_id").
		Where("ar.admin_id = ? AND p.code = ?", adminID, code).
		Count(&c).Error
	if err != nil {
		return false, errors.WithStack(err)
	}
	return c > 0, nil
}

func AssignRoleToAdmin(adminID int64, roleName string) error {
	var role model.Role
	if err := db.Where("name = ?", roleName).First(&role).Error; err != nil {
		return errors.WithStack(err)
	}
	var cnt int64
	db.Model(&model.AdminRole{}).Where("admin_id = ? AND role_id = ?", adminID, role.ID).Count(&cnt)
	if cnt == 0 {
		return db.Create(&model.AdminRole{AdminID: adminID, RoleID: role.ID}).Error
	}
	return nil
}

func CreateRole(name, desc string) error {
	r := &model.Role{Name: name, Desc: desc}
	return errors.WithStack(db.Create(r).Error)
}

func CreatePermission(code, desc string) error {
	p := &model.Permission{Code: code, Desc: desc}
	return errors.WithStack(db.Create(p).Error)
}

func BindPermissionToRole(roleName, permCode string) error {
	var r model.Role
	if err := db.Where("name = ?", roleName).First(&r).Error; err != nil {
		return errors.WithStack(err)
	}
	var p model.Permission
	if err := db.Where("code = ?", permCode).First(&p).Error; err != nil {
		return errors.WithStack(err)
	}
	var cnt int64
	db.Model(&model.RolePermission{}).
		Where("role_id = ? AND permission_id = ?", r.ID, p.ID).Count(&cnt)
	if cnt == 0 {
		return db.Create(&model.RolePermission{RoleID: r.ID, PermissionID: p.ID}).Error
	}
	return nil
}

func UpdateAdminPasswordWithRBAC(id int64, password string) error {
	u := model.Admin{}
	u.ID = id
	if password != "" {
		salt := time.Now().Format("20060102150405")
		u.Password = model.TwoHashPwd(password, salt)
		u.Salt = salt
	}
	u.UpdateTime = time.Now()
	return db.Model(&u).Updates(u).Error
}
