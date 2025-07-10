package verify

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/jordan-wright/email"
	"main/config"
	"net/http"
	"net/smtp"
	"strings"
	"time"
)

var verificationStorage = make(map[string]string)

type SendRequest struct {
	To string `json:"to"`
}

func SendEmailHandler(w http.ResponseWriter, r *http.Request) {
	var req SendRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if req.To == "" {
		http.Error(w, "Email адрес не указан", http.StatusBadRequest)
		return
	}

	cfg := config.Load()

	data := req.To + time.Now().String()
	hash := fmt.Sprintf("%x", md5.Sum([]byte(data)))

	verificationStorage[hash] = req.To

	e := email.NewEmail()
	e.From = cfg.Email
	e.To = []string{req.To}
	e.Subject = "Подтверждение email"
	e.Text = []byte(fmt.Sprintf("Для подтверждения email перейдите по ссылке: http://example.com/verify/%s", hash))
	err := e.Send(cfg.Address, smtp.PlainAuth("", cfg.Email, cfg.Password, extractHost(cfg.Address)))
	if err != nil {
		http.Error(w, "Ошибка отправки email", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"hash":   hash,
	})
}

func VerifyHandler(w http.ResponseWriter, r *http.Request, hash string) {
	email, exists := verificationStorage[hash]

	if !exists {
		http.Error(w, "Неверный код подтверждения", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"email":  email,
	})

	delete(verificationStorage, hash)
}

func extractHost(addr string) string {
	parts := strings.Split(addr, ":")
	if len(parts) > 0 {
		return parts[0]
	}
	return addr
}
