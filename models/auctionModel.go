package models

import (
	"time"
)

type Auction struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Status      string    `json:"status"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	StartPrice  int64     `json:"start_price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// It has the auction details and the auction items
type AuctionResponse struct {
	Auction      Auction       `json:"auction"`
	AuctionItems []AuctionItem `json:"auction_items"`
}

type CreateAuctionRequest struct {
	UserID      int64     `json:"user_id"`
	Status      string    `json:"status"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	StartPrice  int64     `json:"start_price"`
}

type AuctionItem struct {
	ID          int64     `json:"id"`
	AuctionID   int64     `json:"auction_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
type CreateAuctionItemRequest struct {
	AuctionID   int64  `json:"auction_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

type Bid struct {
	ID            int64     `json:"id"`
	AuctionItemID int64     `json:"auction_item_id"`
	UserID        int64     `json:"user_id"`
	Amount        int64     `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
