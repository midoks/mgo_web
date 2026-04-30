package op

import (
	// "fmt"
	// "time"

	// "github.com/pkg/errors"
	// "gorm.io/gorm"

	"mgo/internal/conf"
	"mgo/internal/db"
	// "mgo/internal/model"
	// utils "mgo/internal/utils"
)

func AddLog(uid int64, content string) error {
	if conf.Security.InstallLock {
		return db.AddLog(nil, uid, content)
	}
	return nil

}

func SysLog(content string) error {
	return db.AddLog(nil, 0, content)
}
