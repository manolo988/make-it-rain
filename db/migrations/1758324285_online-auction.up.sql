-- TABLE Auction

CREATE TYPE auction_status AS ENUM ('active', 'inactive', 'completed', 'cancelled');
CREATE TABLE IF NOT EXISTS auction (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    status auction_status NOT NULL default 'active',
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    start_price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE INDEX idx_auction_user_id ON auction(user_id);

-- TABLE AuctionItem
CREATE TABLE IF NOT EXISTS auction_item (
    id BIGSERIAL PRIMARY KEY,
    auction_id BIGINT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    image_url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (auction_id) REFERENCES auction(id)
);

-- TABLE Bid
CREATE TABLE IF NOT EXISTS bid (
    id BIGSERIAL PRIMARY KEY,
    auction_item_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (auction_item_id) REFERENCES auction_item(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE INDEX idx_bid_auction_item_id_created_at ON bid(auction_item_id, created_at);