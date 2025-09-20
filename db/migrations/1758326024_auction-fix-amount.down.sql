-- Revert auctions table
ALTER TABLE auction ALTER COLUMN start_price TYPE DECIMAL(10, 2) USING (start_price::decimal / 100);
ALTER TABLE auction DROP COLUMN currency;
COMMENT ON COLUMN auction.start_price IS NULL;


-- Revert bid table
ALTER TABLE bid ALTER COLUMN amount TYPE DECIMAL(10, 2) USING (amount::decimal / 100);
ALTER TABLE bid DROP COLUMN currency;
COMMENT ON COLUMN bid.amount IS NULL;
