package form

type AdminRecipients struct {
	Name      string `form:"name"`
	MediaType string `form:"media_type"`
}

type AdminRecipientsGroup struct {
	ID     int64  `form:"id"`
	Name   string `form:"name"`
	Status bool   `form:"status"`
}
