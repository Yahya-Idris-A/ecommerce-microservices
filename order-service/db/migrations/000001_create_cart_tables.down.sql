-- Hapus tabel dari yang paling memiliki dependency (anak dulu, baru bapak)
DROP TABLE IF EXISTS cart_items CASCADE;
DROP TABLE IF EXISTS carts CASCADE;