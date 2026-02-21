-- Menambahkan kolom merchant_id dan menjadikannya wajib (NOT NULL)
ALTER TABLE products ADD COLUMN merchant_id UUID NOT NULL;

-- Membuat index untuk mempercepat pencarian produk berdasarkan toko
CREATE INDEX idx_products_merchant ON products(merchant_id);