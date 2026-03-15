-- 1. Hapus tabel user_addresses (Cascading akan menghapus indexnya juga)
DROP TABLE IF EXISTS user_addresses CASCADE;

-- 2. Kembalikan default role menjadi 'customer'
ALTER TABLE users ALTER COLUMN role SET DEFAULT 'customer';

-- 3. Hapus kolom-kolom baru
ALTER TABLE users DROP COLUMN IF EXISTS avatar_url;
ALTER TABLE users DROP COLUMN IF EXISTS phone_number;
ALTER TABLE users DROP COLUMN IF EXISTS full_name;

-- 4. Kembalikan nama kolom password_hash menjadi password
ALTER TABLE users RENAME COLUMN password_hash TO password;