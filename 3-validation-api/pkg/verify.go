package pkg

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/http"
	"strings"
)

type SendRequest struct {
	Email string `json:"email"`
}

func GetScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	return "http"
}

func ExtractHost(addr string) string {
	parts := strings.Split(addr, ":")
	if len(parts) > 0 {
		return parts[0]
	}
	return addr
}

func GenerateToken(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

func ParseAndValidateSendRequest(r *http.Request) (*SendRequest, error) {
	var req SendRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		return nil, fmt.Errorf("некорректный JSON: %w", err)
	}
	if !govalidator.IsEmail(req.Email) {
		return nil, fmt.Errorf("некорректный email")
	}
	return &req, nil
}
