package form

type AdminRecipients struct {
	ID            int64  `form:"id"`
	Name          string `form:"name"`
	MediaType     string `form:"media_type"`
	Mark          string `form:"mark"`
	HashLife      int64  `form:"hash_life"`
	Token         string `form:"token"`
	EmailSmtp     string `form:"email_smtp"`
	EmailUsername string `form:"email_username"`
	EmailPassword string `form:"email_password"`
	EmailFrom     string `form:"email_from"`
	WebhookUrl    string `form:"webhook_url"`
	WebhookMethod string `form:"webhook_method"`
	Status        bool   `form:"status"`
}

type AdminRecipientsGroup struct {
	ID     int64  `form:"id"`
	Name   string `form:"name"`
	Status bool   `form:"status"`
}
