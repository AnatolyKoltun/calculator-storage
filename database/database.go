package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Connect(dsn string) {

	var err error
	DB, err = pgxpool.New(context.Background(), dsn)

	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}

	err = DB.Ping(context.Background())

	if err != nil {
		log.Fatalf("Не удалось выполнить ping БД: %v", err)
	}

	fmt.Println("Подключение к PostgreSQL успешно установлено")

}

func Close() {
	if r := recover(); r != nil {
		log.Fatal(r)
	}

	if DB != nil {
		DB.Close()
	}
}
