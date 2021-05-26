package server

import (
	"log"
	"net/http"

	"github.com/worming004/TwelveFactorApp/mail"
)

type Server struct {
	*http.Server
	mail.MailSender
}

func NewServer(mailSender mail.MailSender, address string, openApiContent []byte) Server {
	routes := getHandlers(mailSender, openApiContent)
	serv := http.Server{
		Addr:    address,
		Handler: routes,
	}
	return Server{&serv, mailSender}
}

func (s *Server) Run() {
	log.Fatal(s.ListenAndServe())
}
