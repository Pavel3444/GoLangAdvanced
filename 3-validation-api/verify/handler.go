package verify

import (
	"encoding/json"
	"fmt"
	"github.com/jordan-wright/email"
	"log"
	"main/config"
	"main/pkg"
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

	file, _ := json.MarshalIndent(data, "", "  ")
	_ = os.WriteFile(verificationFile, file, 0644)
}

func SendEmailHandler(w http.ResponseWriter, r *http.Request) {
	req, err := pkg.ParseAndValidateSendRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cfg := config.Load()
	if cfg == nil {
		log.Fatal("SMTP config is incomplete or missing; application cannot proceed.")
	}

	hash := pkg.GenerateToken(32)

	data := loadVerificationMap()
	data[hash] = req.Email
	saveVerificationMap(data)

	e := email.NewEmail()
	//e.From = cfg.Email
	e.From = "qwe@qwe.qwe"
	e.To = []string{req.Email}
	e.Subject = "Подтверждение email"
	link := fmt.Sprintf("%s://%s/verify/%s", pkg.GetScheme(r), r.Host, hash)
	e.Text = []byte(fmt.Sprintf("Для подтверждения email перейдите по ссылке: %s", link))
	err = e.Send(cfg.Address, smtp.PlainAuth("", cfg.Email, cfg.Password, pkg.ExtractHost(cfg.Address)))
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
	if !exists {
		json.NewEncoder(w).Encode(map[string]bool{
			"verified": false,
		})
		return
	}

	delete(data, hash)
	saveVerificationMap(data)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"verified": true,
		"email":    email,
	})
}
