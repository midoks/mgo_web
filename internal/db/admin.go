package db

import (
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"mgo/internal/model"
	utils "mgo/internal/utils"
)

func GetAdminList(page, size int) ([]model.Admin, int64, error) {
	adminM := db.Model(&model.Admin{})
	var count int64
	if err := adminM.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get server count")
	}

	var list []model.Admin
	if err := db.Order(columnName("id")).Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return list, count, nil
}

func GetAdminById(id int64) (*model.Admin, error) {
	var u model.Admin
	if err := db.First(&u, id).Error; err != nil {
		return nil, errors.Wrapf(err, "failed get admin")
	}
	return &u, nil
}

func GetAdminByName(username string) (*model.Admin, error) {
	info := model.Admin{Username: username}
	if err := db.Where(info).First(&info).Error; err != nil {
		return nil, errors.Wrapf(err, "failed find admin")
	}
	return &info, nil
}

func AdminUpdateEmail(id int64, email string) error {
	return db.Model(&model.Admin{ID: id}).Update("email", email).Error
}

func AdminUpdateTel(id int64, tel string) error {
	return db.Model(&model.Admin{ID: id}).Update("tel", tel).Error
}

func AdminUpdatePass(id int64, password string) error {
	u := model.Admin{}
	u.ID = id

	if password != "" {
		salt := utils.RandString(16)
		u.Password = model.TwoHashPwd(password, salt)
		u.Salt = salt
	}
	u.UpdateTime = time.Now()
	return UpdateAdmin(&u)
}

func UpdateAdmin(u *model.Admin) error {
	if u.Password == "" {
		if err := db.Model(u).Updates(map[string]interface{}{"password": u.Password, "update_time": u.UpdateTime}).Error; err != nil {
			return errors.WithStack(err)
		}
	} else {
		if err := db.Model(u).Updates(u).Error; err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func InitAdmin(user string, pass string) error {
	_, err := GetAdminById(1)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			salt := utils.RandString(16)
			admin := &model.Admin{
				Username: user,
				Password: model.TwoHashPwd(pass, salt),
				Salt:     salt,
			}

			admin.CreateTime = time.Now()
			admin.UpdateTime = time.Now()
			if err := CreateAdmin(admin); err != nil {
				return err
			}
		}
	}
	return nil
}

func CreateAdmin(u *model.Admin) error {
	return errors.WithStack(db.Create(u).Error)
}
