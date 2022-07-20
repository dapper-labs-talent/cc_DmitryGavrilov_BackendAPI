package api

import (
	"encoding/json"
	"net/http"
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
