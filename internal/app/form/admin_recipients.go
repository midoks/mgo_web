package form

type AdminRecipients struct {
	ID        int64  `form:"id"`
	Name      string `form:"name"`
	MediaType string `form:"media_type"`
	Mark      string `form:"mark"`
	HashLife  int64  `form:"hash_life"`
	Status    bool   `form:"status"`
}

type AdminRecipientsGroup struct {
	ID     int64  `form:"id"`
	Name   string `form:"name"`
	Status bool   `form:"status"`
}
