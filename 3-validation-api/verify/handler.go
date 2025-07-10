package verify

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/jordan-wright/email"
	"main/config"
	"net/http"
	"net/smtp"
	"strings"
	"sync"
)

var verificationStorage = sync.Map{}

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

	hash := generateToken(32)

	verificationStorage.Store(hash, req.To)

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
	value, exists := verificationStorage.Load(hash)
	if !exists {
		http.Error(w, "Неверный код подтверждения", http.StatusNotFound)
		return
	}

	email := value.(string)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"email":  email,
	})

	verificationStorage.Delete(hash)
}

func extractHost(addr string) string {
	parts := strings.Split(addr, ":")
	if len(parts) > 0 {
		return parts[0]
	}
	return addr
}

func generateToken(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}
