package services

import (
	"context"

	"github.com/manuel/make-it-rain/db"
	"github.com/manuel/make-it-rain/models"
)

type JobScheduler struct {
	dbService db.DBService
}

func NewJobScheduler(dbService db.DBService) *JobScheduler {
	return &JobScheduler{
		dbService: dbService,
	}
}

func (s *JobScheduler) SubmitJob(
	ctx context.Context,
	req *db.SubmitJobRequest,
) (*models.Job, error) {
	return s.dbService.SubmitJob(ctx, req)
}

func (s *JobScheduler) GetJob(ctx context.Context, userID int64) (*models.Job, error) {
	return s.dbService.GetJob(ctx, userID)
}

func (s *JobScheduler) GetJobs(
	ctx context.Context,
	page, pageSize int,
	sortBy, sortOrder string,
	status string,
) (*db.PaginatedJobs, error) {
	return s.dbService.GetJobs(ctx, page, pageSize, sortBy, sortOrder, status)
}

func (s *JobScheduler) SearchJob(
	ctx context.Context,
	metadata map[string]interface{},
) ([]models.Job, error) {
	return s.dbService.SearchJob(ctx, metadata)
}
