package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/manuel/make-it-rain/db"
	"github.com/manuel/make-it-rain/models"
)

type UserService struct {
	dbService db.DBService
}

func NewUserService(dbService db.DBService) *UserService {
	return &UserService{
		dbService: dbService,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *db.CreateUserRequest) (*models.User, error) {
	req.Password = hashPassword(req.Password)
	return s.dbService.CreateUser(ctx, req)
}

func (s *UserService) GetUser(ctx context.Context, userID int64) (*models.User, error) {
	return s.dbService.GetUser(ctx, userID)
}

func (s *UserService) GetUsers(ctx context.Context, page, pageSize int, sortBy, sortOrder string) (*db.PaginatedUsers, error) {
	return s.dbService.GetUsers(ctx, page, pageSize, sortBy, sortOrder)
}

func (s *UserService) UpdateUser(ctx context.Context, userID int64, updates map[string]interface{}) error {
	return s.dbService.UpdateUser(ctx, userID, updates)
}

func (s *UserService) DeleteUser(ctx context.Context, userID int64) error {
	return s.dbService.DeleteUser(ctx, userID)
}

func (s *UserService) AuthenticateUser(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.dbService.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user.Password != hashPassword(password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("user account is inactive")
	}

	return user, nil
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}