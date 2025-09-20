package db

import (
	"context"
	"fmt"

	"github.com/manuel/make-it-rain/models"
)

type Auction = models.Auction

type PaginatedAuctions struct {
	Auctions   []models.Auction `json:"auctions"`
	TotalCount int              `json:"total_count"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
}

func (s *RealDBService) CreateAuction(
	ctx context.Context,
	auction *models.CreateAuctionRequest,
) (*models.Auction, error) {
	query := `
		INSERT INTO auction (user_id, status, title, description, start_date, end_date, start_price, currency, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING id, user_id, status, title, description, start_date, end_date, start_price, currency, created_at, updated_at`

	var a Auction
	err := Conn.QueryRow(ctx, query, auction.UserID, auction.Status, auction.Title, auction.Description, auction.StartDate, auction.EndDate, auction.StartPrice, auction.Currency).
		Scan(
			&a.ID,
			&a.UserID,
			&a.Status,
			&a.Title,
			&a.Description,
			&a.StartDate,
			&a.EndDate,
			&a.StartPrice,
			&a.Currency,
			&a.CreatedAt,
			&a.UpdatedAt,
		)
	if err != nil {
		return nil, fmt.Errorf("failed to create auction: %w", err)
	}

	return &a, nil
}

func (s *RealDBService) CreateAuctionItem(
	ctx context.Context,
	auctionItem *models.CreateAuctionItemRequest,
) (*models.AuctionItem, error) {
	query := `
		INSERT INTO auction_item (auction_id, name, description, image_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, auction_id, name, description, image_url, created_at, updated_at`

	var ai models.AuctionItem
	err := Conn.QueryRow(ctx, query, auctionItem.AuctionID, auctionItem.Name, auctionItem.Description, auctionItem.ImageURL).
		Scan(&ai.ID, &ai.AuctionID, &ai.Name, &ai.Description, &ai.ImageURL, &ai.CreatedAt, &ai.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create auction item: %w", err)
	}
	return &ai, nil
}
func (s *RealDBService) GetAuctions(
	ctx context.Context,
	page, pageSize int,
	sortBy, sortOrder string,
) (*PaginatedAuctions, error) {
	// Whitelist allowed sort columns to prevent SQL injection
	allowedSortColumns := map[string]bool{
		"id":          true,
		"created_at":  true,
		"updated_at":  true,
		"start_date":  true,
		"end_date":    true,
		"start_price": true,
		"title":       true,
		"status":      true,
	}

	// Validate and sanitize sortBy
	if !allowedSortColumns[sortBy] {
		sortBy = "created_at" // Default to created_at if invalid
	}

	// Validate and sanitize sortOrder
	if sortOrder != "asc" && sortOrder != "ASC" && sortOrder != "desc" && sortOrder != "DESC" {
		sortOrder = "DESC" // Default to DESC if invalid
	}

	countQuery := `SELECT COUNT(*) FROM auction`
	var totalCount int
	err := Conn.QueryRow(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to count auctions: %w", err)
	}
	offset := (page - 1) * pageSize

	// Safely construct query with validated inputs
	query := fmt.Sprintf(`
		SELECT id, user_id, status, title, description, start_date, end_date, start_price, currency, created_at, updated_at
		FROM auction
		ORDER BY %s %s
		LIMIT $1 OFFSET $2`, sortBy, sortOrder)

	rows, err := Conn.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get auctions: %w", err)
	}
	defer rows.Close()

	var auctions []models.Auction
	for rows.Next() {
		var a models.Auction
		err := rows.Scan(
			&a.ID,
			&a.UserID,
			&a.Status,
			&a.Title,
			&a.Description,
			&a.StartDate,
			&a.EndDate,
			&a.StartPrice,
			&a.Currency,
			&a.CreatedAt,
			&a.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan auction: %w", err)
		}
		auctions = append(auctions, a)
	}
	return &PaginatedAuctions{
		Auctions:   auctions,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

// GetAuction by id
func (s *RealDBService) GetAuction(
	ctx context.Context,
	auctionID int64,
) (*models.AuctionResponse, error) {
	query := `SELECT id, user_id, status, title, description, start_date, end_date, start_price, currency, created_at, updated_at FROM auction WHERE id = $1`

	var a models.Auction
	err := Conn.QueryRow(ctx, query, auctionID).Scan(
		&a.ID,
		&a.UserID,
		&a.Status,
		&a.Title,
		&a.Description,
		&a.StartDate,
		&a.EndDate,
		&a.StartPrice,
		&a.Currency,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get auction: %w", err)
	}
	var auctionItems []models.AuctionItem
	query = `SELECT id, auction_id, name, description, image_url, created_at, updated_at FROM auction_item WHERE auction_id = $1`
	rows, err := Conn.Query(ctx, query, auctionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get auction items: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var ai models.AuctionItem
		err := rows.Scan(
			&ai.ID,
			&ai.AuctionID,
			&ai.Name,
			&ai.Description,
			&ai.ImageURL,
			&ai.CreatedAt,
			&ai.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan auction item: %w", err)
		}
		auctionItems = append(auctionItems, ai)
	}
	return &models.AuctionResponse{Auction: a, AuctionItems: auctionItems}, nil
}
