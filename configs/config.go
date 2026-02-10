package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const Port = 3000

type AppConfig struct {
	Db         DbConfig
	Auth       AuthConfig
	Mail       MailConfig
	ServerPort int
}

type DbConfig struct {
	Name string
}

type AuthConfig struct {
	Secret string
}

type MailConfig struct {
	Mail     string
	Password string
	SmtpHost string
	SmtpPort string
}

func GetAppConfig() *AppConfig {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file, using default config.")
	}

	return &AppConfig{
		Db: DbConfig{
			Name: os.Getenv("DB_NAME"),
		},
		Auth: AuthConfig{
			Secret: os.Getenv("TOKEN"),
		},
		ServerPort: Port,
		Mail: MailConfig{
			Mail:     os.Getenv("MAIL"),
			Password: os.Getenv("MAIL_PASS"),
			SmtpHost: os.Getenv("SMTP_HOST"),
			SmtpPort: os.Getenv("SMTP_PORT"),
		},
	}
}
