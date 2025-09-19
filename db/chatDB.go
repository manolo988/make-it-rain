package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/manuel/make-it-rain/models"
)

type Chat = models.Chat

var ErrChatNotFound = errors.New("chat not found")

func (s *RealDBService) CreateChat(
	ctx context.Context,
	chat *models.CreateChatRequest,
) (*Chat, error) {
	query := `
		INSERT INTO chat (owner_id, name, description, chat_type,)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, owner_id, name, description, chat_type, created_at, updated_at`

	var c Chat
	err := Conn.QueryRow(ctx, query,
		chat.OwnerID,
		chat.Name,
		chat.Description,
		chat.ChatType,
	).Scan(
		&c.ID,
		&c.OwnerID,
		&c.Name,
		&c.Description,
		&c.ChatType,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create chat: %w", err)
	}

	for _, userID := range chat.Users {
		_, err := Conn.Exec(
			ctx,
			"INSERT INTO chat_user (chat_id, user_id, created_at, updated_at) VALUES ($1, $2, NOW(), NOW())",
			c.ID,
			userID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to add user to chat: %w", err)
		}
	}

	return &c, nil
}

func (s *RealDBService) GetChat(
	ctx context.Context,
	chatID int64,
) (*Chat, error) {
	query := `
		SELECT id, owner_id, name, description, chat_type, created_at, updated_at
		FROM chat
		WHERE id = $1`

	var c Chat
	err := Conn.QueryRow(ctx, query, chatID).Scan(
		&c.ID,
		&c.OwnerID,
		&c.Name,
		&c.Description,
		&c.ChatType,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrChatNotFound
		}
		return nil, fmt.Errorf("failed to get chat: %w", err)
	}

	owner, err := s.GetUser(ctx, c.OwnerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat owner: %w", err)
	}
	c.Owner = *owner

	usersQuery := `
		SELECT u.id, u.email, u.name, u.password, u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN chat_user cu ON cu.user_id = u.id
		WHERE cu.chat_id = $1`
	rows, err := Conn.Query(ctx, usersQuery, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.Name,
			&u.Password,
			&u.IsActive,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan chat user: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate chat users: %w", err)
	}
	c.Users = users

	messagesQuery := `
		SELECT id, chat_id, content, sender_id, created_at, updated_at
		FROM message
		WHERE chat_id = $1
		ORDER BY created_at`
	msgRows, err := Conn.Query(ctx, messagesQuery, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat messages: %w", err)
	}
	defer msgRows.Close()

	var messages []models.Message
	for msgRows.Next() {
		var m models.Message
		if err := msgRows.Scan(
			&m.ID,
			&m.ChatID,
			&m.Content,
			&m.SenderID,
			&m.CreatedAt,
			&m.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan chat message: %w", err)
		}
		messages = append(messages, m)
	}
	if err := msgRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate chat messages: %w", err)
	}
	c.Messages = messages

	return &c, nil
}

func (s *RealDBService) SendMessage(
	ctx context.Context,
	message *models.CreateMessageRequest,
) (*models.Message, error) {
	query := `
		INSERT INTO message (chat_id, content, sender_id, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, chat_id, content, sender_id, created_at, updated_at`

	var m models.Message
	err := Conn.QueryRow(ctx, query,
		message.ChatID,
		message.Content,
		message.SenderID,
	).Scan(
		&m.ID,
		&m.ChatID,
		&m.Content,
		&m.SenderID,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	return &m, nil
}
