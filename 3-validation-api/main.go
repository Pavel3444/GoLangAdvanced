package main

import (
	"github.com/joho/godotenv"
	"log"
	"main/verify"
	"net/http"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	verify.SetupRoutes()

	log.Println("Сервер запущен на :8081")
	err = http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
