DROP INDEX IF EXISTS idx_categories_merchant_id;
ALTER TABLE categories DROP COLUMN IF EXISTS merchant_id;