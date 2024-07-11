package server

import (
	"errors"
	"fmt"
	"gohttp/utils"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func HandleRouter() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			SendHTTPErrorResponse(res, http.StatusMethodNotAllowed)
			return
		}

		if len(req.UserAgent()) == 0 {
			SendHTTPErrorResponse(res, http.StatusForbidden)
			return
		}

		Router(res, req)
	}
}

func Router(res http.ResponseWriter, req *http.Request) {
	parsedURL, err := url.Parse(req.URL.Path)
	if err != nil {
		SendHTTPErrorResponse(res, http.StatusBadRequest)
		return
	}

	path := parsedURL.Path
	switch path {
	case "/":
		fullPath := filepath.Join("html", "index.html")
		SendStaticFile(res, req, fullPath)
	default:
		fullPath := filepath.Join("html", filepath.Clean(path))
		fmt.Println(fullPath)
		SendStaticFile(res, req, fullPath)
	}
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
