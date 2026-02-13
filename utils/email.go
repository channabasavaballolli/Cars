package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
)

// SendOTP sends the verification code via email using SMTP
// If SMTP credentials are not set, it logs the code to the console (Dev Mode)
func SendOTP(email, code string) error {
	smtpHost := strings.TrimSpace(os.Getenv("SMTP_HOST"))
	smtpPort := strings.TrimSpace(os.Getenv("SMTP_PORT"))
	smtpEmail := strings.TrimSpace(os.Getenv("SMTP_EMAIL"))
	smtpPassword := strings.TrimSpace(os.Getenv("SMTP_PASSWORD"))

	// Console Mode (Dev)
	if smtpHost == "" || smtpPort == "" || smtpEmail == "" || smtpPassword == "" {
		log.Printf("--------------------------------------------------")
		log.Printf("[DEV MODE] Email Service Skipped (Missing .env vars)")
		log.Printf("To: %s", email)
		log.Printf("Subject: Your Login Code")
		log.Printf("Body: Your verification code is: %s", code)
		log.Printf("--------------------------------------------------")
		return nil
	}

	// Real Email Sending
	auth := smtp.PlainAuth("", smtpEmail, smtpPassword, smtpHost)
	to := []string{email}
	msg := []byte("To: " + email + "\r\n" +
		"Subject: Your Login Code\r\n" +
		"\r\n" +
		"Your verification code is: " + code + "\r\n")

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	err := smtp.SendMail(addr, auth, smtpEmail, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
