
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
	message.SetHeader("Subject", "Reset Your Password - goAuth")

	htmlBody := fmt.Sprintf(`
	<div style="font-family: Arial, sans-serif; max-width: 500px; margin: auto; padding: 25px; text-align: center; border: 1px solid #e5e7eb; border-radius: 10px;">
		<h2 style="color: #2563eb;">Reset Your Password</h2>

		<p style="color: #374151; line-height: 1.6;">
			We received a request to reset your goAuth account password.
		</p>

		<div style="background-color: #f3f4f6; padding: 16px; margin: 20px 0; border-radius: 6px; overflow-wrap: anywhere;">
			<p style="font-size: 12px; color: #6b7280;">
				YOUR PASSWORD RESET TOKEN
			</p>
			<p style="font-family: monospace; font-size: 16px; font-weight: bold; color: #111827;">
				%s
			</p>
		</div>

		<p style="font-size: 14px; color: #4b5563;">
			Enter this token in the password reset form to choose a new password.
		</p>

		<p style="margin-top: 25px; font-size: 12px; color: #9ca3af;">
			If you did not request a password reset, ignore this email.
		</p>
	</div>
	`, token)

	message.SetBody("text/html", htmlBody)

	dialer := gomail.NewDialer(
		"smtp.gmail.com",
		587,
		sender,
		appPassword,
	)

	if err := dialer.DialAndSend(message); err != nil {
		return false, fmt.Errorf("failed to send password reset email: %w", err)
	}

	fmt.Println("Password reset email sent successfully to", receiver)

	return true, nil
}
