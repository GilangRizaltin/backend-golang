package pkg

import (
	"net/smtp"
	"os"
	"strconv"
)

// const (
// 	smtpServer  = "smtp.gmail.com"
// 	smtpPort    = 587
// 	senderEmail = os.Getenv("SENDER_EMAIL")
// 	senderPass  = os.Getenv("SENDER_PASS")
// )

func SendEmail(to []string, subject, body string) error {
	smtpServer := "smtp.gmail.com"
	smtpPort := 587
	senderEmail := os.Getenv("SENDER_EMAIL")
	senderPass := os.Getenv("SENDER_PASS")
	auth := smtp.PlainAuth("", senderEmail, senderPass, smtpServer)
	msg := []byte("To: " + to[0] + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")
	err := smtp.SendMail(smtpServer+":"+strconv.Itoa(smtpPort), auth, senderEmail, to, msg)
	if err != nil {
		return err
	}
	return nil
}
