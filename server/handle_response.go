package server

import (
	"errors"
	"gohttp/utils"
	"io"
	"net/http"
	"os"
)

func SendHTTPErrorResponse(res http.ResponseWriter, status int) {
	msg := "gohttp: " + http.StatusText(status)
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(status)
	io.WriteString(res, msg)
}

func SendStaticFile(res http.ResponseWriter, req *http.Request, filepath string) {
	isDir, err := utils.VerifyPath(filepath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			SendHTTPErrorResponse(res, http.StatusNotFound)
		} else if os.IsPermission(err) {
			SendHTTPErrorResponse(res, http.StatusForbidden)
		} else {
			SendHTTPErrorResponse(res, http.StatusInternalServerError)
		}
		return
	} else if isDir {
		SendHTTPErrorResponse(res, http.StatusNotFound)
		return
	}

	http.ServeFile(res, req, filepath)
}
