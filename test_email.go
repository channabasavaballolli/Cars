//go:build ignore

package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	smtpHost := strings.TrimSpace(os.Getenv("SMTP_HOST"))
	smtpPort := strings.TrimSpace(os.Getenv("SMTP_PORT"))
	smtpEmail := strings.TrimSpace(os.Getenv("SMTP_EMAIL"))
	smtpPassword := strings.TrimSpace(os.Getenv("SMTP_PASSWORD"))

	fmt.Printf("SMTP Configuration:\n")
	fmt.Printf("Host: %s\n", smtpHost)
	fmt.Printf("Port: %s\n", smtpPort)
	fmt.Printf("Email: %s\n", smtpEmail)
	fmt.Printf("Password Length: %d\n", len(smtpPassword))

	if len(smtpPassword) != 16 {
		fmt.Println("WARNING: Gmail App Passwords are typically 16 characters long.")
		fmt.Println("Yours is", len(smtpPassword), "characters.")
		fmt.Println("Did you include spaces? We are trimming them automatically in this test.")
	}

	// Message
	to := []string{smtpEmail} // Send to self
	msg := []byte("To: " + smtpEmail + "\r\n" +
		"Subject: Test Email from Go\r\n" +
		"\r\n" +
		"This is a test email to verify SMTP settings.\r\n")

	// Auth
	auth := smtp.PlainAuth("", smtpEmail, smtpPassword, smtpHost)

	// Send
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	fmt.Printf("\nAttempting to send email to %s...\n", smtpEmail)

	err = smtp.SendMail(addr, auth, smtpEmail, to, msg)
	if err != nil {
		fmt.Printf("\nERROR: %v\n", err)
		fmt.Println("\nTroubleshooting Tips:")
		fmt.Println("1. Ensure 'SMTP_EMAIL' matches the account that generated the App Password.")
		fmt.Println("2. Ensure 'SMTP_PASSWORD' is the 16-character App Password (no spaces).")
		fmt.Println("3. Ensure 2-Step Verification is ON for your Google Account.")
	} else {
		fmt.Println("\nSUCCESS! Email sent successfully.")
		fmt.Println("This confirms your .env credentials are correct.")
	}
}
