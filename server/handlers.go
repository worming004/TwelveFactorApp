package server

import (
	_ "embed"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/worming004/TwelveFactorApp/mail"
)

func getHandlers(sender mail.MailSender, openApiContent []byte) *mux.Router {
	router := mux.NewRouter()
	router.Handle("/mail", postMailHandler(sender)).Methods("POST")
	router.HandleFunc("/openapi.yaml", serveOpenApi(openApiContent)).Methods("GET")
	router.HandleFunc("/openapi.yml", serveOpenApi(openApiContent)).Methods("GET")

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

func postMailHandler(sender mail.MailSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request PostMailRequest
		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		err = sender.SendMail(mail.Mail{
			To:      request.To,
			Subject: request.Subject,
			Body:    []byte(request.Body),
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusCreated)
	}
}
