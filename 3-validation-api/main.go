package main

import (
	"log"
	"main/verify"
	"net/http"
)

func main() {
	verify.SetupRoutes()

	log.Println("Сервер запущен на :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
