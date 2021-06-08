package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/worming004/TwelveFactorApp/auth"
	"github.com/worming004/TwelveFactorApp/log"
	"github.com/worming004/TwelveFactorApp/mail"
)

func getHandlers(sender mail.MailSender, openApiContent []byte, eventRepository EventRepository, jwtWrap auth.JwtWrapper, version string) *mux.Router {
	router := mux.NewRouter()
	authMiddleware := getAuthMiddleware(jwtWrap)
	router.Handle("/mail", authMiddleware(postMailHandler(sender, eventRepository))).Methods("POST")
	router.HandleFunc("/openapi.yaml", serveOpenApi(openApiContent)).Methods("GET")
	router.HandleFunc("/openapi.yml", serveOpenApi(openApiContent)).Methods("GET")
	router.HandleFunc("/jwt", jwtWrap.CreateToken).Methods("POST")
	var versionHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("/version called")
		w.Write([]byte(version))
	}
	router.HandleFunc("/version", versionHandler).Methods("GET")
	router.Handle("/auth/version", authMiddleware(versionHandler)).Methods("GET")
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}).Methods("GET")

	return router
}

func serveOpenApi(openApiContent []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(openApiContent)
	}
}

type PostMailRequest struct {
	To      string `json:"To"`
	Subject string `json:"Subject"`
	Body    string `json:"Body"`
}

type Event struct {
	Subject string
}
type EventRepository interface {
	CreateEvent(Event) error
}

func postMailHandler(sender mail.MailSender, eventRepository EventRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request PostMailRequest
		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			log.WriteErrorLog(w, err.Error(), http.StatusBadRequest)
		}

		m := mail.Mail{
			To:      request.To,
			Subject: request.Subject,
			Body:    []byte(request.Body),
		}

		err = sender.SendMail(m)

		if err != nil {
			log.WriteErrorLog(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = eventRepository.CreateEvent(Event{m.Subject})

		if err != nil {
			log.WriteErrorLog(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		logrus.Info("send mail ok")
		w.WriteHeader(http.StatusCreated)
	}
}
