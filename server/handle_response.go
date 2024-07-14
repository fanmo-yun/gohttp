package server

import (
	"errors"
	"gohttp/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func SendHTTPErrorResponse(response http.ResponseWriter, status int) {
	msg := "gohttp: " + http.StatusText(status)
	response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	response.WriteHeader(status)
	io.WriteString(response, msg)
}

func SendStaticFile(response http.ResponseWriter, request *http.Request, path string) {
	handleError := func(err error) {
		if errors.Is(err, os.ErrNotExist) {
			SendHTTPErrorResponse(response, http.StatusNotFound)
		} else if os.IsPermission(err) {
			SendHTTPErrorResponse(response, http.StatusForbidden)
		} else {
			SendHTTPErrorResponse(response, http.StatusInternalServerError)
		}
	}

	isDir, err := utils.VerifyPath(path)
	if err != nil {
		handleError(err)
		return
	}

	if isDir {
		indexPath := filepath.Join(path, "index.html")
		_, err := utils.VerifyPath(indexPath)
		if err != nil {
			handleError(err)
			return
		}
		path = indexPath
	}

	http.ServeFile(response, request, path)
}
