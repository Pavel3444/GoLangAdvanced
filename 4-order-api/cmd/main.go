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

	if err := migrations.InitialMigration(dataBase); err != nil {
		fmt.Printf("❌ Failed to run migrations: %v\n", err)
		return
	}
	fmt.Println("✅ Database connection established")
}
