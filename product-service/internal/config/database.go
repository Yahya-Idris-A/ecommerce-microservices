package config

import (
	"database/sql"
	"fmt"
	"log"

	// Blank identifier (_) digunakan di sini.
	// Ini konsep penting di Golang: kita tidak memanggil fungsi dari package ini secara langsung,
	// tapi kita butuh fungsi init() di dalamnya berjalan untuk mendaftarkan driver 'postgres' ke database/sql.
	_ "github.com/lib/pq"
)

// InitDB membuka koneksi ke database PostgreSQL
func InitDB(host, port, user, password, dbname string) *sql.DB {
	// DSN (Data Source Name) adalah format string koneksi standar
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Membuka pool koneksi ke database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open database connection: %v", err)
	}

	// sql.Open sebenarnya tidak langsung melakukan koneksi.
	// Kita butuh Ping() untuk memastikan kredensialnya benar dan database menyala.
	if err = db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	log.Println("successfully connected to the database")
	return db
}
