package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/manuel/make-it-rain/models"
)

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type User = models.User

type PaginatedUsers struct {
	Users      []models.User `json:"users"`
	TotalCount int           `json:"total_count"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
}

func (s *RealDBService) CreateUser(ctx context.Context, user *CreateUserRequest) (*User, error) {
	query := `
		INSERT INTO users (email, name, password, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, true, NOW(), NOW())
		RETURNING id, email, name, password, is_active, created_at, updated_at`

	var u User
	err := Conn.QueryRow(ctx, query,
		user.Email,
		user.Name,
		user.Password,
	).Scan(
		&u.ID,
		&u.Email,
		&u.Name,
		&u.Password,
		&u.IsActive,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &u, nil
}

func (s *RealDBService) GetUser(ctx context.Context, userID int64) (*User, error) {
	query := `
		SELECT id, email, name, password, is_active, created_at, updated_at
		FROM users
		WHERE id = $1`

	var u User
	err := Conn.QueryRow(ctx, query, userID).Scan(
		&u.ID,
		&u.Email,
		&u.Name,
		&u.Password,
		&u.IsActive,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &u, nil
}

func (s *RealDBService) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, name, password, is_active, created_at, updated_at
		FROM users
		WHERE email = $1`

	var u User
	err := Conn.QueryRow(ctx, query, email).Scan(
		&u.ID,
		&u.Email,
		&u.Name,
		&u.Password,
		&u.IsActive,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &u, nil
}

func (s *RealDBService) GetUsers(ctx context.Context, page, pageSize int, sortBy, sortOrder string) (*PaginatedUsers, error) {
	countQuery := `SELECT COUNT(*) FROM users`
	var totalCount int
	err := Conn.QueryRow(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT id, email, name, password, is_active, created_at, updated_at
		FROM users
		ORDER BY %s %s
		LIMIT $1 OFFSET $2`, sortBy, sortOrder)

	rows, err := Conn.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.Name,
			&u.Password,
			&u.IsActive,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, u)
	}

	return &PaginatedUsers{
		Users:      users,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s *RealDBService) UpdateUser(ctx context.Context, userID int64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	setClauses := []string{}
	args := []interface{}{userID}
	argCount := 1

	for field, value := range updates {
		argCount++
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, argCount))
		args = append(args, value)
	}

	query := fmt.Sprintf(`
		UPDATE users
		SET %s, updated_at = NOW()
		WHERE id = $1`,
		joinStrings(setClauses, ", "))

	result, err := Conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (s *RealDBService) DeleteUser(ctx context.Context, userID int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := Conn.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}