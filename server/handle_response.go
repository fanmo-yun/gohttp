package server

import (
	"io"
	"net/http"
)

func SendHTTPErrorResponse(res http.ResponseWriter, status int) {
	msg := "gohttp: " + http.StatusText(status)
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(status)
	io.WriteString(res, msg)
}
