package config

import "os"

type Config struct {
	SMTPAddress string
	Email       string
	Password    string
}

func Load() *Config {
	return &Config{
		SMTPAddress: os.Getenv("SMTP_ADDRESS"),
		Email:       os.Getenv("SMTP_EMAIL"),
		Password:    os.Getenv("SMTP_PASSWORD"),
	}
}
