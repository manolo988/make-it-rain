package services

import (
	"context"

	"github.com/manuel/make-it-rain/db"
	"github.com/manuel/make-it-rain/models"
)

type UserMessageService struct {
	dbService db.DBService
}

func NewUserMessageService(dbService db.DBService) *UserMessageService {
	return &UserMessageService{
		dbService: dbService,
	}
}

func (s *UserMessageService) GetUserMessages(
	ctx context.Context,
	userID int64,
) ([]models.Message, error) {
	return s.dbService.GetUserMessages(ctx, userID)
}
