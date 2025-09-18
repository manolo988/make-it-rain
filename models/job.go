package models

import (
	"encoding/json"
	"time"
)

type Job struct {
	ID            int64           `json:"id"`
	UserID        int64           `json:"user_id"`
	ScheduledTime time.Time       `json:"scheduled_time"`
	Metadata      json.RawMessage `json:"metadata"`
	PartitionID   string          `json:"partition_id"`
	JobStatus     string          `json:"job_status"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}
