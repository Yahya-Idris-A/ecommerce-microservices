-- 1. Ubah nama kolom password menjadi password_hash
ALTER TABLE users RENAME COLUMN password TO password_hash;

-- 2. Tambahkan kolom baru ke tabel users
-- Catatan: Kita kasih DEFAULT '' sementara agar tidak error jika sudah ada data lama di tabel
ALTER TABLE users ADD COLUMN full_name VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN phone_number VARCHAR(50) UNIQUE;
ALTER TABLE users ADD COLUMN avatar_url TEXT;

-- 3. Ubah default role dari 'customer' menjadi 'buyer'
ALTER TABLE users ALTER COLUMN role SET DEFAULT 'buyer';

-- 4. Buat tabel user_addresses baru
CREATE TABLE IF NOT EXISTS user_addresses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    label VARCHAR(100) NOT NULL,
    recipient_name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(50) NOT NULL,
    full_address TEXT NOT NULL,
    city VARCHAR(100) NOT NULL,
    province VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20) NOT NULL,
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexing untuk mempercepat tarikan data alamat saat checkout
CREATE INDEX IF NOT EXISTS idx_user_addresses_user_id ON user_addresses(user_id);