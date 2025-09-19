package services

import (
	"context"

	"github.com/manuel/make-it-rain/db"
	"github.com/manuel/make-it-rain/models"
)

type ChatService struct {
	dbService db.DBService
}

func NewChatService(dbService db.DBService) *ChatService {
	return &ChatService{
		dbService: dbService,
	}
}

func (s *ChatService) CreateChat(
	ctx context.Context,
	req *models.CreateChatRequest,
) (*models.Chat, error) {
	return s.dbService.CreateChat(ctx, req)
}

func (s *ChatService) SendMessage(
	ctx context.Context,
	req *models.CreateMessageRequest,
) (*models.Message, error) {
	return s.dbService.SendMessage(ctx, req)
}

func (s *ChatService) GetChat(
	ctx context.Context,
	chatID int64,
) (*models.Chat, error) {
	return s.dbService.GetChat(ctx, chatID)
}
