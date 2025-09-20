-- Fix auctions table
ALTER TABLE auction ADD COLUMN currency VARCHAR(3) NOT NULL DEFAULT 'USD';
UPDATE auction SET start_price = start_price * 100;
ALTER TABLE auction ALTER COLUMN start_price TYPE BIGINT;
COMMENT ON COLUMN auction.start_price IS 'in cents';


-- Fix bid table
ALTER TABLE bid ADD COLUMN currency VARCHAR(3) NOT NULL DEFAULT 'USD';
UPDATE bid SET amount = amount * 100;
ALTER TABLE bid ALTER COLUMN amount TYPE BIGINT;
COMMENT ON COLUMN bid.amount IS 'in cents'; 