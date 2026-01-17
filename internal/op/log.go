package op

import (
	// "fmt"
	// "time"

	// "github.com/pkg/errors"
	// "gorm.io/gorm"

	"mgo/internal/db"
	// "mgo/internal/model"
	// utils "mgo/internal/utils"
)

func AddLog(uid int64, content string) error {
	return db.AddLog(uid, content)
}
