package model

import (
	"encoding/json"
	"time"
)

type ClusterNodeLogin struct {
	ID         int64     `json:"id" gorm:"primaryKey"`  // unique key
	Name       string    `json:"name"`                  // name
	NodeID     int64     `json:"node_id" gorm:"unique"` // node_id
	Params     string    `json:"ip"`                    // params
	Status     bool      `json:"status"`                // status
	CreateTime time.Time `json:"create_time"`           // create_time
	UpdateTime time.Time `json:"update_time"`           // update_time
}

type ClusterNodeLoginParams struct {
	Host  string `json:"host"`
	Port  string `json:"port"`
	SshID string `json:"ssh_id"`
}

func (a *ClusterNodeLogin) SetParams(p ClusterNodeLoginParams) error {
	b, err := json.Marshal(p)
	if err != nil {
		return err
	}
	a.Params = string(b)
	return nil
}

func (a *ClusterNodeLogin) GetParams() (ClusterNodeLoginParams, error) {
	var p ClusterNodeLoginParams
	if a.Params == "" {
		return p, nil
	}
	err := json.Unmarshal([]byte(a.Params), &p)
	return p, err
}
