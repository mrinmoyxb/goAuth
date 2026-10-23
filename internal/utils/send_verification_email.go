package utils

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

func SendVerificationEmail(receiver string, token string) (bool, error) {
	message := gomail.NewMessage()
	sender := os.Getenv("SENDER_MAIL")
	appPassword := os.Getenv("MAIL_PASSWORD")
	mailPORT := os.Getenv("PORT")
	if sender == "" || appPassword == "" {
		return false, errors.New("invalid mail credentials")
	}

	message.SetHeader("From", sender)
	message.SetHeader("To", receiver)
	message.SetHeader("Subject", "Verify You Email - goAuth")
	fmt.Printf("sender: %s, receiver: %s", sender, receiver)

	verificationURL := fmt.Sprintf("http://localhost:%s/goAuth/verify-email?token=%s&email=%s", mailPORT, token, receiver)

	htmlBody := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; text-align: center; padding: 20px;">
			<h2>Hello!</h2>
			<p>Thank you for signing up. Please click the button below to verify your email address:</p>
			<a href="%s" style="background-color: #4CAF50; color: white; padding: 12px 20px; text-decoration: none; border-radius: 5px; display: inline-block; font-weight: bold; margin-top: 15px;">
				Verify Email
			</a>
			<p style="margin-top: 20px; font-size: 12px; color: #888;">If you didn't request this, you can safely ignore this email.</p>
		</div>
	`, verificationURL)

	message.SetBody("text/html", htmlBody)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, sender, appPassword)
	if err := dialer.DialAndSend(message); err != nil {
		return false, errors.New("error sending mails")
	}

	return true, nil
}
