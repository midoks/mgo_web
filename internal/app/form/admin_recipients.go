package form

type AdminRecipients struct {
	ID            int64  `form:"id"`
	Name          string `form:"name" binding:"required"`
	MediaType     string `form:"media_type" binding:"required"`
	Mark          string `form:"mark"`
	HashLife      int64  `form:"hash_life"`
	Token         string `form:"token"`
	SendID        string `form:"send_id"`
	EmailSmtp     string `form:"email_smtp"`
	EmailUsername string `form:"email_username"`
	EmailPassword string `form:"email_password"`
	EmailFrom     string `form:"email_from"`
	WebhookUrl    string `form:"webhook_url"`
	WebhookMethod string `form:"webhook_method"`
	Count         int64  `form:"count" binding:"required"`
	Minutes       int64  `form:"minutes" binding:"required"`
	Status        bool   `form:"status"`
}

type AdminRecipientsGroup struct {
	ID     int64  `form:"id"`
	Name   string `form:"name"`
	Status bool   `form:"status"`
}
