package form

type AdminRecipients struct {
	Name      string `form:"name"`
	MediaType string `form:"media_type"`
}
