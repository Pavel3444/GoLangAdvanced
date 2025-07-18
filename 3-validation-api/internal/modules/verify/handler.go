package verify

import (
	"3-validation-api/config"
	"3-validation-api/pkg/verify"
	"encoding/json"
	"fmt"
	"github.com/jordan-wright/email"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"sync"
)

var mu sync.Mutex

const verificationFile = "verification.json"

func loadVerificationMap() map[string]string {
	mu.Lock()
	defer mu.Unlock()

	data := make(map[string]string)
	file, err := os.ReadFile(verificationFile)
	if err != nil {
		return data
	}
	_ = json.Unmarshal(file, &data)
	return data
}

func saveVerificationMap(data map[string]string) {
	mu.Lock()
	defer mu.Unlock()

	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("ошибка сериализации verification map: %v", err)
		return
	}

	err = os.WriteFile(verificationFile, file, 0644)
	if err != nil {
		log.Printf("ошибка записи verification файла: %v", err)
	}
}

func SendEmailHandler(w http.ResponseWriter, r *http.Request) {
	req, err := verify.ParseAndValidateSendRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cfg := config.Load()
	if cfg == nil {
		log.Printf("SMTP конфигурация отсутствует - функционал отправки email недоступен")
		http.Error(w, "Сервис отправки email временно недоступен", http.StatusServiceUnavailable)
		return

	}

	hash := verify.GenerateToken(32)

	data := loadVerificationMap()
	data[hash] = req.Email
	saveVerificationMap(data)

	e := email.NewEmail()
	e.From = cfg.Email
	e.To = []string{req.Email}
	e.Subject = "Подтверждение email"
	link := fmt.Sprintf("%s://%s/verify/%s", verify.GetScheme(r), r.Host, hash)
	e.Text = []byte(fmt.Sprintf("Для подтверждения email перейдите по ссылке: %s", link))
	err = e.Send(cfg.Address, smtp.PlainAuth("", cfg.Email, cfg.Password, verify.ExtractHost(cfg.Address)))
	if err != nil {
		log.Printf("ошибка отправки письма: %v", err)
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
	w.Header().Set("Content-Type", "application/json")

	data := loadVerificationMap()
	email, exists := data[hash]
	if exists {
		delete(data, hash)
		saveVerificationMap(data)
	}

	if !exists {
		json.NewEncoder(w).Encode(map[string]bool{
			"verified": false,
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"verified": true,
		"email":    email,
	})
}
