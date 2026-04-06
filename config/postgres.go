package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type DataSourceName struct {
	DatabaseURL string
}

func (dsn *DataSourceName) GetDatabaseURL() {
	//if err := godotenv.Load(); err != nil {
	//	panic("Файл .env не найден")
	//}

	// Пробуем загрузить .env только если файл существует (для локальной разработки)
	if _, err := os.Stat(".env"); err == nil {
		// Файл существует, загружаем
		_ = godotenv.Load() // Игнорируем ошибку
	}

	dsn.DatabaseURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	fmt.Println(dsn.DatabaseURL)
}
