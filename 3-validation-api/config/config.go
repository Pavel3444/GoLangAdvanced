package config

import (
	"os"
	"strings"
)

type Config struct {
	Address  string
	Email    string
	Password string
}

func Load() *Config {
	address := ensureSMTPPort(os.Getenv("SMTP_ADDRESS"))
	email := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")

	if address == "" || email == "" || password == "" {
		return nil
	}

	return &Config{
		Address:  address,
		Email:    email,
		Password: password,
	}
}

func ensureSMTPPort(addr string) string {
	if addr == "" {
		return ""
	}
	if !strings.Contains(addr, ":") {
		return addr + ":587"
	}
	return addr
}
