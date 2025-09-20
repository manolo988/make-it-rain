package services

import (
	"context"

	"github.com/manuel/make-it-rain/db"
	"github.com/manuel/make-it-rain/models"
)

type AuctionService struct {
	dbService db.DBService
}

func NewAuctionService(dbService db.DBService) *AuctionService {
	return &AuctionService{
		dbService: dbService,
	}
}

func (s *AuctionService) CreateAuction(
	ctx context.Context,
	req *models.CreateAuctionRequest,
) (*models.Auction, error) {
	return s.dbService.CreateAuction(ctx, req)
}

func (s *AuctionService) CreateAuctionItem(
	ctx context.Context,
	req *models.CreateAuctionItemRequest,
) (*models.AuctionItem, error) {
	return s.dbService.CreateAuctionItem(ctx, req)
}

func (s *AuctionService) GetAuctions(
	ctx context.Context,
	page, pageSize int,
	sortBy, sortOrder string,
) (*db.PaginatedAuctions, error) {
	return s.dbService.GetAuctions(ctx, page, pageSize, sortBy, sortOrder)
}

func (s *AuctionService) GetAuction(
	ctx context.Context,
	auctionID int64,
) (*models.AuctionResponse, error) {
	return s.dbService.GetAuction(ctx, auctionID)
}
