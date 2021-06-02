package main

import (
	"context"
	_ "embed"
	"os"
	"time"

	"github.com/worming004/TwelveFactorApp/auth"
	"github.com/worming004/TwelveFactorApp/db"
	"github.com/worming004/TwelveFactorApp/mail"
	"github.com/worming004/TwelveFactorApp/server"
)

//go:embed openapi.yaml
var openApiContent []byte

type nullMailSender struct{}

func (n nullMailSender) SendMail(m mail.Mail) error {
	return nil
}

func main() {
	mongoClient, err := db.NewDefaultMongoClient()
	if err != nil {
		panic(err)
	}

	contextTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer mongoClient.Disconnect(contextTimeout)
	eventRepository := db.NewMongoEventRepository(mongoClient)
	address := os.Getenv("TWELVE_SERVER_ADDRESS")
	sender := getMailSender()
	jwtWrap := auth.GetDefaultJwtWrapper()
	server := server.NewServer(sender, address, openApiContent, eventRepository, jwtWrap)

	server.Run()
}

func getMailSender() mail.MailSender {
	conf := mail.GetConfigurationByEnvironmentVariable()
	sender := mail.NewMailSender(conf)
	// sender := nullMailSender{}

	return sender
}
