package server

import (
	"errors"
	"gohttp/utils"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func Handle(res http.ResponseWriter, req *http.Request) *url.URL {
	if req.Method != http.MethodGet {
		SendHTTPErrorResponse(res, http.StatusMethodNotAllowed)
		return nil
	}

	if len(req.UserAgent()) == 0 {
		SendHTTPErrorResponse(res, http.StatusForbidden)
		return nil
	}

	parsedURL, err := url.Parse(req.URL.Path)
	if err != nil {
		SendHTTPErrorResponse(res, http.StatusBadRequest)
		return nil
	}

	return parsedURL
}

func HandleRouter(config *utils.Config) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		urlPath := Handle(res, req)
		if urlPath == nil {
			return
		}

		path := urlPath.Path
		handled := CustomRouter(res, req, path, config.Static.Dirpath, config.Custom)
		if !handled {
			Router(res, req, path, config.Static)
		}
	}
}

func CustomRouter(res http.ResponseWriter, req *http.Request, URLPath string, staticDir string, cus []utils.CustomConfig) bool {
	for _, custom := range cus {
		if custom.Urlpath == URLPath {
			fullPath := filepath.Join(staticDir, filepath.Clean(custom.Filepath))
			SendStaticFile(res, req, fullPath)
			return true
		}
	}
	return false
}

func Router(res http.ResponseWriter, req *http.Request, URLPath string, h utils.HtmlConfig) {
	switch URLPath {
	case "/":
		fullPath := filepath.Join(h.Dirpath, h.Index)
		SendStaticFile(res, req, fullPath)
	default:
		fullPath := filepath.Join(h.Dirpath, filepath.Clean(URLPath))
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
