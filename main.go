package main

import (
	_ "embed"
	"os"

	"github.com/worming004/TwelveFactorApp/mail"
	"github.com/worming004/TwelveFactorApp/server"
)

//go:embed openapi.yaml
var openApiContent []byte

func main() {
	conf := mail.GetConfigurationByEnvironmentVariable()
	address := os.Getenv("TWELVE_SERVER_ADDRESS")
	sender := mail.NewMailSender(conf)
	server := server.NewServer(sender, address, openApiContent)

	server.Run()
}
