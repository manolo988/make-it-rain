-- Rollback for online auction feature
-- Drop indexes first where applicable, then tables in dependency order, then enum type

-- Bid
DROP INDEX IF EXISTS idx_bid_auction_item_id_created_at;
DROP TABLE IF EXISTS bid;

-- AuctionItem
DROP TABLE IF EXISTS auction_item;

-- Auction
DROP INDEX IF EXISTS idx_auction_user_id;
DROP TABLE IF EXISTS auction;

-- Enum type
DROP TYPE IF EXISTS auction_status;

