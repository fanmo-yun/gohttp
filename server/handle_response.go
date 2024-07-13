package server

import (
	"errors"
	"gohttp/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func SendHTTPErrorResponse(res http.ResponseWriter, status int) {
	msg := "gohttp: " + http.StatusText(status)
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(status)
	io.WriteString(res, msg)
}

func SendStaticFile(res http.ResponseWriter, req *http.Request, path string) {
	handleError := func(err error) {
		if errors.Is(err, os.ErrNotExist) {
			SendHTTPErrorResponse(res, http.StatusNotFound)
		} else if os.IsPermission(err) {
			SendHTTPErrorResponse(res, http.StatusForbidden)
		} else {
			SendHTTPErrorResponse(res, http.StatusInternalServerError)
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

	http.ServeFile(res, req, path)
}
