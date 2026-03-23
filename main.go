package main

import (
	"github.com/AnatolyKoltun/calculator-storage/config"
	"github.com/AnatolyKoltun/calculator-storage/database"
)

func connectToDB() {
	defer database.Close()

	dsn := new(config.DataSourceName)
	dsn.GetDatabaseURL()

	database.Connect(dsn.DatabaseURL)
}

func main() {
	connectToDB()
}
