package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/manuel/make-it-rain/models"
)

type SubmitJobRequest struct {
	UserID        int64           `json:"user_id"        binding:"required"`
	ScheduledTime time.Time       `json:"scheduled_time" binding:"required"`
	Metadata      json.RawMessage `json:"metadata"       binding:"required"`
	PartitionID   string          `json:"partition_id"   binding:"required"`
}

type Job = models.Job

type PaginatedJobs struct {
	Jobs       []models.Job `json:"jobs"`
	TotalCount int          `json:"total_count"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
}

// CREATE TABLE IF NOT EXISTS job (
//     id BIGSERIAL PRIMARY KEY,
//     user_id BIGINT NOT NULL,
//     scheduled_time TIMESTAMP NOT NULL,
//     metadata JSONB,
//     partition_id TEXT,
//     job_status status_enum DEFAULT 'submitted',
//     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
//     updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
//     FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
// );

func (s *RealDBService) SubmitJob(ctx context.Context, job *SubmitJobRequest) (*Job, error) {
	query := `
		INSERT INTO job (user_id, scheduled_time, metadata, partition_id, job_status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, 'submitted', NOW(), NOW())
		RETURNING id, user_id, scheduled_time, metadata, partition_id, job_status, created_at, updated_at`

	var j Job
	err := Conn.QueryRow(ctx, query,
		job.UserID,
		job.ScheduledTime,
		job.Metadata,
		job.PartitionID,
	).Scan(
		&j.ID,
		&j.UserID,
		&j.ScheduledTime,
		&j.Metadata,
		&j.PartitionID,
		&j.JobStatus,
		&j.CreatedAt,
		&j.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	return &j, nil
}

func (s *RealDBService) GetJob(ctx context.Context, jobID int64) (*Job, error) {
	query := `
		SELECT id, user_id, scheduled_time, metadata, partition_id, job_status, created_at, updated_at
		FROM job
		WHERE id = $1`

	var j Job
	err := Conn.QueryRow(ctx, query, jobID).Scan(
		&j.ID,
		&j.UserID,
		&j.ScheduledTime,
		&j.Metadata,
		&j.PartitionID,
		&j.JobStatus,
		&j.CreatedAt,
		&j.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("job not found")
		}
		return nil, fmt.Errorf("failed to get job: %w", err)
	}
	return &j, nil
}

func (s *RealDBService) GetJobs(
	ctx context.Context,
	page, pageSize int,
	sortBy, sortOrder, status string,
) (*PaginatedJobs, error) {
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortOrder == "" {
		sortOrder = "desc"
	}

	allowedSortBy := map[string]bool{
		"id":             true,
		"user_id":        true,
		"scheduled_time": true,
		"partition_id":   true,
		"job_status":     true,
		"created_at":     true,
		"updated_at":     true,
	}
	if !allowedSortBy[sortBy] {
		return nil, fmt.Errorf("invalid sort_by parameter")
	}

	allowedSortOrder := map[string]bool{
		"asc":  true,
		"desc": true,
	}
	if !allowedSortOrder[sortOrder] {
		return nil, fmt.Errorf("invalid sort_order parameter")
	}

	countQuery := `SELECT COUNT(*) FROM job WHERE job_status = $1`
	var totalCount int
	err := Conn.QueryRow(ctx, countQuery, status).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to count jobs: %w", err)
	}

	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT id, user_id, scheduled_time, metadata, partition_id, job_status, created_at, updated_at
		FROM job
		WHERE job_status = $3
		ORDER BY %s %s
		LIMIT $1 OFFSET $2
	`, sortBy, sortOrder)

	rows, err := Conn.Query(ctx, query, pageSize, offset, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs: %w", err)
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var j Job
		err := rows.Scan(
			&j.ID,
			&j.UserID,
			&j.ScheduledTime,
			&j.Metadata,
			&j.PartitionID,
			&j.JobStatus,
			&j.CreatedAt,
			&j.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan job: %w", err)
		}
		jobs = append(jobs, j)
	}

	return &PaginatedJobs{
		Jobs:       jobs,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s *RealDBService) SearchJob(
	ctx context.Context,
	metadata map[string]interface{},
) ([]Job, error) {
	query := `
		SELECT id, user_id, scheduled_time, metadata, partition_id, job_status, created_at, updated_at
		FROM job
		WHERE metadata @> $1`

	rows, err := Conn.Query(ctx, query, metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to get job: %w", err)
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var j Job
		err := rows.Scan(
			&j.ID,
			&j.UserID,
			&j.ScheduledTime,
			&j.Metadata,
			&j.PartitionID,
			&j.JobStatus,
			&j.CreatedAt,
			&j.UpdatedAt,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("job not found")
			}
			return nil, fmt.Errorf("failed to get job: %w", err)
		}
		jobs = append(jobs, j)
	}

	return jobs, nil
}
