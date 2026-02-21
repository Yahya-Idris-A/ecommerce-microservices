-- Hapus produk dulu karena dia bergantung pada kategori
DROP TABLE IF EXISTS products;

-- Setelah produk bersih, baru kita bisa menghapus kategori dengan aman
DROP TABLE IF EXISTS categories;