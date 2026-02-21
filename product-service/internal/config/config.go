package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config adalah struktur data untuk menampung semua konfigurasi aplikasi
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	AppPort    string
}

// LoadConfig membaca file .env dan mengembalikannya dalam bentuk struct
func LoadConfig() *Config {
	// Memuat file .env yang ada di root direktori
	err := godotenv.Load()
	if err != nil {
		// Kita pakai log.Println, bukan log.Fatal.
		// Kenapa? Karena saat di-deploy ke production (misal Kubernetes Docker),
		// kita tidak pakai file .env, melainkan langsung inject ke OS environment.
		log.Println("Warning: .env file not found, reading from OS environment")
	}

	return &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		AppPort:    os.Getenv("PORT"),
	}
}
