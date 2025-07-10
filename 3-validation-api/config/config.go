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
	return &Config{
		Address:  ensureSMTPPort(os.Getenv("SMTP_ADDRESS")),
		Email:    os.Getenv("SMTP_EMAIL"),
		Password: os.Getenv("SMTP_PASSWORD"),
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
