package api

import (
	"context"
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
				writePlainTextError(jsonError, w, r)
			}
		}
	default:
		{
			writePlainTextError(err, w, r)
		}
	}
}

func writePlainTextError(err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		logrus.Error(errors.Wrap(err, "internal server error"))
	}
	w.WriteHeader(http.StatusInternalServerError)
	content := []byte(`{"code":500, "errorMessage":"internal server error"}`)
	_, err = w.Write(content)
	if err != nil {
		logrus.Error(errors.Wrap(err, "error writing information to a http response"))
	}
}

type contextMiddlewareHandler func(w http.ResponseWriter, r *http.Request) (context.Context, error)

func (imh contextMiddlewareHandler) execute(w http.ResponseWriter, r *http.Request, next http.Handler) {
	ctx, err := imh(w, r)
	if err != nil {
		unwrapErrorResponse(err, w, r)
		return
	}

	if ctx != nil {
		r = r.WithContext(ctx)
	}

	next.ServeHTTP(w, r)

}

func (imh contextMiddlewareHandler) handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		imh.execute(w, r, next)
	})
}

func identityMiddleware(imh contextMiddlewareHandler) func(http.Handler) http.Handler {
	return imh.handler
}
