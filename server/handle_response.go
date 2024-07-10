package server

import (
	"io"
	"net/http"
)

func SendHTTPErrorResponse(w http.ResponseWriter, status int) {
	msg := "gohttp: " + http.StatusText(status)
	w.WriteHeader(http.StatusMethodNotAllowed)
	io.WriteString(w, msg)
}
