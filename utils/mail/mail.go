package mail

import (
	"fmt"
	"log/slog"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func Send(confirmationCode, to string) {
	if err := godotenv.Load(); err != nil {
		slog.Error("error reading environment variables")
	}
	from, exists := os.LookupEnv("EMAIL")
	if !exists {
		slog.Error("Can't read EMAIL")
		return
	}
	pass, exists := os.LookupEnv("EMAIL_PASSWORD")
	if !exists {
		slog.Error("Can't read EMAIL_PASSWORD")
		return
	}
	msg := generateMessage(confirmationCode, to, from)
	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		slog.Error("smtp error: %s", err)
		return
	}
	slog.Info("Successfully sended to " + to)
}
func generateMessage(confirmationCode, to, from string) string {
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Sport bot email confirmation" + "\n\n" +
		fmt.Sprintf("click on link: http://127.0.0.1:9000/users/confirm-email/%v to confirm your email", confirmationCode)
	return msg
}
