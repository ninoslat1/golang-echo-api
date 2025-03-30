package utils

import (
	models "echo-api/models"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"gopkg.in/gomail.v2"
)

func SendVerificationEmail(toEmail, securityCode string) error {
	err := godotenv.Load()
	if err != nil {
		log.Warn("Error loading .env file, using default environment variables")
	}

	cfg := &models.EmailConfig{
		Password: os.Getenv("SMTP_PASSWORD"),
		Email:    os.Getenv("HOST_EMAIL"),
	}

	fromEmail := cfg.Email
	password := cfg.Password

	smtpHost := "smtp.gmail.com"
	smtpPort := 587

	message := gomail.NewMessage()
	message.SetHeader("From", fromEmail)
	message.SetHeader("To", toEmail)
	message.SetHeader("Subject", "Verifikasi Akun Anda")
	message.SetBody("text/html", fmt.Sprintf("<h3>Kode Verifikasi Anda: %s</h3>", securityCode))

	dialer := gomail.NewDialer(smtpHost, smtpPort, fromEmail, password)

	if err := dialer.DialAndSend(message); err != nil {
		message := fmt.Sprintf("Failed to send verification email to: %s", toEmail)
		log.Info(message, err)
		return err
	}

	log.Info("Email verification sent successfully to ", toEmail)
	return nil
}
