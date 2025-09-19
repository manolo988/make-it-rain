package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/manuel/make-it-rain/models"
)

type Message = models.Message

var ErrMessageNotFound = errors.New("message not found")

func (s *RealDBService) GetUserMessages(ctx context.Context, userID int64) ([]Message, error) {
	// We need to get all the messages from the user where are undelivered
	// First get all chats of the user
	queryChats := `
		SELECT chat_id
		FROM chat_user
		WHERE user_id = $1
	`
	var chatsIDs []int64
	rows, err := Conn.Query(ctx, queryChats, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user chats: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var c Chat
		if err := rows.Scan(&c.ID); err != nil {
			return nil, fmt.Errorf("failed to scan user chat: %w", err)
		}
		chatsIDs = append(chatsIDs, c.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate user chats: %w", err)
	}
	// Then get all messages from the chats filtered by userID and that are undelivered sorted by created_at. Return the messages
	queryMessages := `
		SELECT id, chat_id, content, sender_id, created_at, updated_at
		FROM message	
		WHERE chat_id = ANY($1)
		ORDER BY created_at`
	var messages []Message
	rows, err = Conn.Query(ctx, queryMessages, chatsIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get user messages: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.ChatID, &m.Content, &m.SenderID, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user message: %w", err)
		}
		messages = append(messages, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate user messages: %w", err)
	}
	return messages, nil
}
