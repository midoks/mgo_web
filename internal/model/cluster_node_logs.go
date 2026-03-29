package model

type ClusterNodeLogs struct {
	ID          int64  `json:"id" gorm:"primaryKey"`                    // unique key
	Day         int64  `json:"day" gorm:"index"`                        // day
	Description string `json:"description"`                             // description
	NodeID      int64  `json:"node_id" gorm:"index" binding:"required"` // node_id
	Level       string `json:"level"`                                   // level
	Tag         string `json:"tag"`                                     // tag
	NodeTime    int64  `json:"node_time"`                               // node_time
	CreateTime  int64  `json:"create_time"`                             // create_time
}
