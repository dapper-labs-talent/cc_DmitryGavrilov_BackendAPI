package api

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	HeaderContentType = "Content-Type"
	ContentTypeJSON   = "application/json"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set(HeaderContentType, ContentTypeJSON)
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.WriteHeader(status)
	_, err = w.Write(bytes)
	return err
}

type identityHttpHandler func(w http.ResponseWriter, r *http.Request) error

func handler(f identityHttpHandler) http.HandlerFunc {
	return f.execute
}

func (ihh identityHttpHandler) execute(w http.ResponseWriter, r *http.Request) {
	err := ihh(w, r)
	if err != nil {
		unwrapErrorResponse(err, w, r)
	}
}

func unwrapErrorResponse(err error, w http.ResponseWriter, r *http.Request) {
	switch e := err.(type) {
	case *HttpErrorResponse:
		{
			jsonError := writeJSON(w, e.Code, e)
			if jsonError != nil {
				unwrapErrorResponse(jsonError, w, r)
			}
		}
	default:
		{
			w.WriteHeader(http.StatusInternalServerError)
			errorData := []byte(`{
				"code":500, "errorMessage":"Internal server error"
			}`)
			_, err := w.Write(errorData)
			if err != nil {
				logrus.Error(errors.Wrap(err, "error writing information to a http response"))
			}
		}
	}
}
