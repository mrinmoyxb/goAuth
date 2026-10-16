package utils

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

func SendMail(receiver string, token string) (bool, error) {
	m := gomail.NewMessage()
	sender := os.Getenv("SENDER_MAIL")
	app_password := os.Getenv("MAIL_PASSWORD")
	if sender == "" || app_password == "" {
		return false, errors.New("invalid mail credentials")
	}

	m.SetHeader("From", sender)
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", "Verify You Email - goAuth")

	verificationURL := fmt.Sprintf("http://localhost:8080/goAuth/verify-email?token=%s?email=%s", token, receiver)

	htmlBody := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; text-align: center; padding: 20px;">
			<h2>Hello!</h2>
			<p>Thank you for signing up. Please click the button below to verify your email address:</p>
			<a href="%s" style="background-color: #4CAF50; color: white; padding: 12px 20px; text-decoration: none; border-radius: 5px; display: inline-block; font-weight: bold; margin-top: 15px;">
				Verify Email
			</a>
			<p style="margin-top: 20px; font-size: 12px; color: #888;">If you didn't request this, you can safely ignore this email.</p>
		</div>"smtp.gmail.com
	`, verificationURL)

	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer("smtp.gmail.com", 587, sender, app_password)
	if err := d.DialAndSend(m); err != nil {
		return false, errors.New("error sending mails")
	}

	fmt.Println("Email sent successfully to", receiver)
	return true, nil
}
