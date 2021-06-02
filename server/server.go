package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/worming004/TwelveFactorApp/auth"
	"github.com/worming004/TwelveFactorApp/mail"
)

type Server struct {
	*http.Server
	mail.MailSender
}

func NewServer(mailSender mail.MailSender, address string, openApiContent []byte, eventRepository EventRepository, jwtWrap auth.JwtWrapper, version string) Server {
	routes := getHandlers(mailSender, openApiContent, eventRepository, jwtWrap, version)
	serv := http.Server{
		Addr:    address,
		Handler: routes,
	}
	return Server{&serv, mailSender}
}

func (s *Server) Run() {
	closingChan := make(chan struct{})
	go func() {
		osSignal := make(chan os.Signal)
		// os signal, ctrl+c
		signal.Notify(osSignal, os.Interrupt)
		// kubernetes kill signal
		signal.Notify(osSignal, syscall.SIGTERM)
		signal.Notify(osSignal, syscall.SIGINT)
		<-osSignal
		if err := s.Shutdown(context.Background()); err != nil {
			logrus.Error("shutdown error: ", err)
		}
		logrus.Info("server is shut down")
		close(closingChan)
	}()
	logrus.Info("server is started on port ", s.Addr)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		logrus.Warn("server not properly closed", err)
	}
	logrus.Info("server is closed")
	<-closingChan
}
