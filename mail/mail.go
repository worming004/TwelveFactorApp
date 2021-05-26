package mail

import (
	"fmt"
	"net/smtp"
	"os"
)

type From string
type Password string

type Configuration struct {
	From     From
	Password Password
}

func GetConfigurationByEnvironmentVariable() Configuration {
	from := os.Getenv("TWELVE_MAIL_FROM")
	pass := os.Getenv("TWELVE_MAIL_PASSWORD")

	return Configuration{From(from), Password(pass)}
}

type Mail struct {
	To      string
	Subject string
	Body    []byte
}

func NewMailSender(c Configuration) MailSender {
	return concreteMailSender{c}
}

type MailSender interface {
	SendMail(Mail) error
}

type concreteMailSender struct {
	Configuration
}

func (ms concreteMailSender) SendMail(m Mail) error {
	// Sender data.
	from := string(ms.From)
	password := string(ms.Password)

	// Receiver email address.
	to := []string{
		m.To,
	}

	// smtp server configuration.
	smtpHost := "send.one.com"
	smtpPort := "587"

	// Message.
	message := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", from, m.To, m.Subject, m.Body)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(message))
	if err != nil {
		return err
	}
	return nil
}
