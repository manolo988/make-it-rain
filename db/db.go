package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Conn *pgxpool.Pool

type DBService interface {
	CreateUser(ctx context.Context, user *CreateUserRequest) (*User, error)
	GetUser(ctx context.Context, userID int64) (*User, error)
	GetUsers(
		ctx context.Context,
		page, pageSize int,
		sortBy, sortOrder string,
	) (*PaginatedUsers, error)
	UpdateUser(ctx context.Context, userID int64, updates map[string]interface{}) error
	DeleteUser(ctx context.Context, userID int64) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	// Job
	SubmitJob(ctx context.Context, job *SubmitJobRequest) (*Job, error)
	GetJob(ctx context.Context, jobID int64) (*Job, error)
	GetJobs(
		ctx context.Context,
		page, pageSize int,
		sortBy, sortOrder string,
		status string,
	) (*PaginatedJobs, error)
	SearchJob(ctx context.Context, metadata map[string]interface{}) ([]Job, error)
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

type RealDBService struct{}

func NewDBService() DBService {
	return &RealDBService{}
}

func InitDB(connStr string) error {
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("failed to parse database config: %w", err)
	}

	config.MaxConns = 20
	config.MinConns = 2
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	Conn = pool
	return nil
}

func CloseDB() {
	if Conn != nil {
		Conn.Close()
	}
}

func (s *RealDBService) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return Conn.Begin(ctx)
}
