-- 1. Aktifkan UUID extension (wajib untuk service baru)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 2. Buat tabel Payung Keranjang (carts)
CREATE TABLE IF NOT EXISTS carts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID UNIQUE NOT NULL, -- 1 User = 1 Cart
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk mempercepat pencarian keranjang berdasarkan user
CREATE INDEX IF NOT EXISTS idx_carts_user_id ON carts(user_id);

-- 3. Buat tabel Isi Keranjang (cart_items)
CREATE TABLE IF NOT EXISTS cart_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cart_id UUID NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    product_id UUID NOT NULL,
    merchant_id UUID NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0), -- Quantity nggak boleh minus/nol
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    -- Pastikan 1 produk di dalam 1 keranjang tidak ada data kembar (duplikat baris)
    UNIQUE(cart_id, product_id)
);

-- Indexing krusial untuk fitur Grouping by Merchant dan Load Cart
CREATE INDEX IF NOT EXISTS idx_cart_items_cart_id ON cart_items(cart_id);
CREATE INDEX IF NOT EXISTS idx_cart_items_merchant_id ON cart_items(merchant_id);