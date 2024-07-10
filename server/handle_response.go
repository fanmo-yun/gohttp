package server

import (
	"io"
	"net/http"
)

func SendHTTPErrorResponse(w http.ResponseWriter, status int) {
	msg := "gohttp: " + http.StatusText(status)
	w.WriteHeader(status)
	io.WriteString(w, msg)
}
