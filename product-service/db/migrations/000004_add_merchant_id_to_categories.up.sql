-- Tambahkan kolom merchant_id yang bersifat opsional (NULLABLE)
ALTER TABLE categories ADD COLUMN merchant_id UUID NULL;

-- Tambahkan index agar pencarian kategori milik suatu toko menjadi cepat
CREATE INDEX idx_categories_merchant_id ON categories(merchant_id);