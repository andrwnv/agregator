package services

import (
	"fmt"
	"net/smtp"
)

type Mailer struct {
	smtpHost string
	smtpPort string
	username string
	password string

	auth smtp.Auth
}

func MakeMailer(host string, port string, password string, username string) *Mailer {
	return &Mailer{
		smtpHost: host,
		smtpPort: port,
		password: password,
		username: username,
		auth:     smtp.PlainAuth("", username, password, host),
	}
}

func (mailer *Mailer) SendVerifyEmail(to []string, uuid string) (err error) {
	from := "no-reply@take-place.com"

	msg := []byte(fmt.Sprintf("From:%s \r\n", from) +
		fmt.Sprintf("To:%s \r\n", to) +
		"Subject: Verify account\r\n\r\n" +
		fmt.Sprintf("Verify link: 127.0.0.1:9090/api/v1/user/verify/%s \r\n", uuid))

	auth := smtp.PlainAuth("", mailer.username, mailer.password, mailer.smtpHost)
	return smtp.SendMail(fmt.Sprintf("%s:%s", mailer.smtpHost, mailer.smtpPort), auth, from, to, msg)
}
