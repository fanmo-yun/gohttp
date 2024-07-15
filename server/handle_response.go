package server

import (
	"errors"
	"gohttp/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
)

func SendHTTPErrorResponse(response http.ResponseWriter, status int) {
	msg := "gohttp: " + http.StatusText(status)
	zap.L().Error("HTTP error response", zap.Int("status", status), zap.String("message", msg))
	response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	response.WriteHeader(status)
	io.WriteString(response, msg)
}

func SendStaticFile(response http.ResponseWriter, request *http.Request, path string) {
	handleError := func(err error) {
		if errors.Is(err, os.ErrNotExist) {
			zap.L().Warn("File not found", zap.String("path", path))
			SendHTTPErrorResponse(response, http.StatusNotFound)
		} else if os.IsPermission(err) {
			zap.L().Warn("Permission denied", zap.String("path", path))
			SendHTTPErrorResponse(response, http.StatusForbidden)
		} else {
			zap.L().Error("Internal server error", zap.Error(err))
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

	if strings.HasSuffix(path, ".ico") {
		response.Header().Set("Content-Type", "image/x-icon")
		http.ServeFile(response, request, path)
		return
	}
	http.ServeFile(response, request, path)
}

func SendTryRootFile(response http.ResponseWriter, request *http.Request, path string, h utils.HtmlConfig) {
	zap.L().Info("Serving file with try-root strategy", zap.String("path", path), zap.String("request_url", request.URL.String()), zap.String("dirpath", h.Dirpath), zap.String("index", h.Index))
	handleError := func(err error) {
		if errors.Is(err, os.ErrNotExist) {
			fullpath := filepath.Join(h.Dirpath, h.Index)
			zap.L().Warn("File not found, trying root file", zap.String("path", path), zap.String("fallback_path", fullpath))
			http.ServeFile(response, request, fullpath)
		} else if os.IsPermission(err) {
			zap.L().Warn("Permission denied", zap.String("path", path))
			SendHTTPErrorResponse(response, http.StatusForbidden)
		} else {
			zap.L().Error("Internal server error", zap.Error(err))
			SendHTTPErrorResponse(response, http.StatusInternalServerError)
		}
	}

	_, err := utils.VerifyPath(path)
	if err != nil {
		handleError(err)
		return
	}

	if strings.HasSuffix(path, ".ico") {
		response.Header().Set("Content-Type", "image/x-icon")
		http.ServeFile(response, request, path)
		return
	}
	http.ServeFile(response, request, path)
}
