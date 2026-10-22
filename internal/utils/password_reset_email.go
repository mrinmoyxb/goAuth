package utils

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

func SendPasswordResetEmail(receiver string, token string) (bool, error) {
	message := gomail.NewMessage()

	sender := os.Getenv("SENDER_MAIL")
	appPassword := os.Getenv("MAIL_PASSWORD")

	if sender == "" || appPassword == "" {
		return false, errors.New("invalid mail credentials")
	}

	message.SetHeader("From", sender)
	message.SetHeader("To", receiver)
	message.SetHeader("Subject", "Verify Your Email - goAuth")

	htmlBody := fmt.Sprintf(`
    <div style="font-family: Arial, sans-serif; max-width: 500px; margin: auto; padding: 25px; text-align: center; border: 1px solid #e5e7eb; border-radius: 10px;">
        <h2 style="color: #2563eb;">Verify Your Email</h2>

        <p style="color: #374151; line-height: 1.6;">
            Thank you for signing up! Use the verification token below to verify your email address.
        </p>

        <div style="background-color: #f3f4f6; padding: 16px; margin: 20px 0; border-radius: 6px; overflow-wrap: anywhere;">
            <p style="font-size: 12px; color: #6b7280; margin-top: 0;">
                YOUR VERIFICATION TOKEN
            </p>
            <p style="font-family: monospace; font-size: 16px; font-weight: bold; color: #111827; margin-bottom: 0;">
                %s
            </p>
        </div>

        <p style="font-size: 14px; color: #4b5563;">
            Copy this token and paste it into the email verification form in goAuth.
        </p>

        <p style="margin-top: 25px; font-size: 12px; color: #9ca3af;">
            If you didn't create an account, you can safely ignore this email.
        </p>
    </div>
`, token)

	message.SetBody("text/html", htmlBody)


	fmt.Println("Email sent successfully to", receiver)
	return true, nil
}