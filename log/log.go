package log

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func WriteErrorLog(w http.ResponseWriter, err string, code int) {
	logrus.WithField("http-status", code).Error(err)
	http.Error(w, err, code)
}
