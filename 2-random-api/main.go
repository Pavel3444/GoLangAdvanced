package main

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/random", func(w http.ResponseWriter, r *http.Request) {
		responseNum := randomInt(1, 6)
		w.Header().Set("Content-Type", "text/plain")
		_, err := w.Write([]byte(strconv.Itoa(responseNum)))
		if err != nil {
			log.Printf("Ошибка при отправке ответа: %v", err)
			return
		}
	})
	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}
	log.Printf("Сервер запущен на http://localhost%s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}
