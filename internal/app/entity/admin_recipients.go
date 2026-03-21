package entity

import (
	"mgo/internal/model"
)

type AdminRecipientsEntityList struct {
	model.AdminRecipients
	Name string `json:"name"` // name
}
