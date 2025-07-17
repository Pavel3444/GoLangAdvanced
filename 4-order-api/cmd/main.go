package main

import (
	"4-order-api/config"
	"4-order-api/migrations"
	"4-order-api/pkg/db"
	"fmt"
)

func main() {
	conf := config.LoadConfig()
	dataBase := db.NewDb(*conf)

	migrations.InitialMigration()
	fmt.Println(dataBase)
}
