package server

import (
	"log"
	"net/http"

	"github.com/worming004/TwelveFactorApp/auth"
	"github.com/worming004/TwelveFactorApp/mail"
)

type Server struct {
	*http.Server
	mail.MailSender
}

func NewServer(mailSender mail.MailSender, address string, openApiContent []byte, eventRepository EventRepository, jwtWrap auth.JwtWrapper) Server {
	routes := getHandlers(mailSender, openApiContent, eventRepository, jwtWrap)
	serv := http.Server{
		Addr:    address,
		Handler: routes,
	}
	return Server{&serv, mailSender}
}

func (s *Server) Run() {
	log.Fatal(s.ListenAndServe())
}
