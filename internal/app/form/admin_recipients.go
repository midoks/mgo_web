package form

type AdminRecipients struct {
	ID          int64  `form:"id" json:"id"`
	AdminID     int64  `form:"admin_id" json:"admin_id"`
	MediaID     int64  `form:"media_id" json:"media_id"`
	RecipientID string `form:"recipients_id" json:"recipient_id"`
	ClusterID   int64  `form:"cluster_id" json:"cluster_id"`
	Mark        string `form:"mark" json:"mark"`
	Status      bool   `form:"status" json:"status"`
}
