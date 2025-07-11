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
	"sync"
)

var verificationStorage = sync.Map{}

func SendEmailHandler(w http.ResponseWriter, r *http.Request) {
	req, err := pkg.ParseAndValidateSendRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cfg := config.Load()
	hash := pkg.GenerateToken(32)
	verificationStorage.Store(hash, req.Email)
	e := email.NewEmail()
	//e.From = cfg.Email
	//TODO: example for test
	e.From = "from@example.com"
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

	value, exists := verificationStorage.Load(hash)
	if !exists {
		json.NewEncoder(w).Encode(map[string]bool{
			"verified": false,
		})
		return
	}

	// Удаляем, как требует ТЗ
	verificationStorage.Delete(hash)

	// Возвращаем JSON с verified: true и email (если нужно)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"verified": true,
		"email":    value,
	})
}
