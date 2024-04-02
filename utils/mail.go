package utils

import (
	"log/slog"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func Send(body string) {
	from := "shiningsuffer@gmail.com"
	if err := godotenv.Load(); err != nil {
		slog.Error("error reading environment variables")
	}
	pass, exists := os.LookupEnv("EMAIL_PASSWORD")
	if !exists {
		slog.Error("Can't read EMAIL_PASSWORD")
		return
	}
	to := "anldeboshir@gmail.com"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Hello there\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		slog.Error("smtp error: %s", err)
		return
	}
	slog.Info("Successfully sended to " + to)
}
