package model

import (
	"encoding/json"
)

type ClusterNode struct {
	ID          int64  `json:"id" gorm:"primaryKey"`                             // unique key
	ClusterID   int64  `json:"cluster_id" gorm:"index"`                          // cluster_id
	Name        string `json:"name"`                                             // name
	IsInstalled bool   `json:"is_installed"`                                     // is_installed
	NodeInfo    string `json:"node_info"`                                        // node_info
	Secret      string `json:"secret" gorm:"unique;index" binding:"required"`    // secret
	UniqueID    string `json:"unique_id" gorm:"unique;index" binding:"required"` // unique_id
	Mark        string `json:"mark"`                                             // mark
	Status      bool   `json:"status" gorm:"index"`                              // status
	CreateTime  int64  `json:"create_time"`                                      // create_time
	UpdateTime  int64  `json:"update_time"`                                      // update_time
}

type ClusterNodeNodeInfoParam struct {
	BuildVersion string `json:"build_version"`
	OS           string `json:"os"`
	Arch         string `json:"arch"`
	ExePath      string `json:"exe_path"`
	Hostname     string `json:"hostname"`
	UpdatedAt    int64  `json:"updated_at"`
	Timestamp    int64  `json:"timestamp"`

	//cpu
	CPUUsage         float64 `json:"cpu_usage"`
	CPULogicalCount  int     `json:"cpu_logical_count"`
	CPUPhysicalCount int     `json:"cpu_physical_count"`

	// mem
	MemoryUsage float64 `json:"memory_usage"`
	MemoryTotal uint64  `json:"memory_total"`

	// load
	Load1m  float64 `json:"load1m"`
	Load5m  float64 `json:"load5m"`
	Load15m float64 `json:"load15m"`

	// disk
	DiskUsage             float64 `json:"disk_usage"`
	DiskMaxUsage          float64 `json:"disk_max_usage"`
	DiskMaxUsagePartition string  `json:"disk_max_usage_partition"`
	DiskTotal             uint64  `json:"disk_total"`
	DiskWritingSpeedMB    int     `json:"disk_writing_speed_mb"` // 硬盘写入速度

	// traffic
	TrafficInBytes  uint64 `json:"traffic_in_bytes"`
	TrafficOutBytes uint64 `json:"traffic_out_bytes"`

	IsActive bool   `json:"is_active"`
	Error    string `json:"error"`
}

func (a *ClusterNode) SetNodeInfoParams(p ClusterNodeNodeInfoParam) error {
	b, err := json.Marshal(p)
	if err != nil {
		return err
	}
	a.NodeInfo = string(b)
	return nil
}

func (a *ClusterNode) GetNodeInfoParams() (ClusterNodeNodeInfoParam, error) {
	var p ClusterNodeNodeInfoParam
	if a.NodeInfo == "" {
		return p, nil
	}
	err := json.Unmarshal([]byte(a.NodeInfo), &p)
	return p, err
}
