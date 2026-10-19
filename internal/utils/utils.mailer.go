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

	verificationURL := fmt.Sprintf("http://localhost:3000/goAuth/verify-email?token=%s&email=%s", token, receiver)

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

func SendMailForgotPassword(receiver string, token string) (bool, error) {
	m := gomail.NewMessage()

	sender := os.Getenv("SENDER_MAIL")
	app_password := os.Getenv("MAIL_PASSWORD")

	if sender == "" || app_password == "" {
		return false, errors.New("invalid mail credentials")
	}

	m.SetHeader("From", sender)
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", "Verify Your Email - goAuth")

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

	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer("smtp.gmail.com", 587, sender, app_password)

	if err := d.DialAndSend(m); err != nil {
		return false, errors.New("error sending mail")
	}

	fmt.Println("Email sent successfully to", receiver)
	return true, nil
}
