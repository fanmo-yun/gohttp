package server

import (
	"gohttpd/utils"
	"net/http"
	"net/url"
	"path/filepath"

	"go.uber.org/zap"
)

func Handle(response http.ResponseWriter, request *http.Request) *url.URL {
	if request.Method != http.MethodGet {
		SendHTTPErrorResponse(response, http.StatusMethodNotAllowed)
		return nil
	}

	if len(request.UserAgent()) == 0 {
		SendHTTPErrorResponse(response, http.StatusForbidden)
		return nil
	}

	parsedURL, err := url.Parse(request.URL.Path)
	if err != nil {
		SendHTTPErrorResponse(response, http.StatusBadRequest)
		return nil
	}

	return parsedURL
}

func HandleRouter(config *utils.Config) http.HandlerFunc {
	proxies := CreateProxies(config.Proxy)
	backends := NewBackend(config.Backend)
	loadBalancer := NewLoadBalancer(backends)

	return func(response http.ResponseWriter, request *http.Request) {
		zap.L().Info("Request received", zap.String("method", request.Method), zap.String("url", request.URL.String()))

		urlPath := Handle(response, request)
		if urlPath == nil {
			zap.L().Warn("Invalid URL path")
			return
		}
		path := urlPath.Path

		if proxies != nil && FindAndServeProxy(response, request, path, proxies) {
			zap.L().Info("Handled by proxy", zap.String("path", path))
			return
		}

		if backends != nil {
			loadBalancer.ServeHTTP(response, request)
			zap.L().Info("Handled by load balancer", zap.String("path", path))
			return
		}

		if !config.Static.Try {
			if !CustomRouter(response, request, path, config.Static.Dirpath, config.Custom, config.Gzip) {
				zap.L().Info("Handled by static router", zap.String("path", path))
				Router(response, request, path, config.Static, config.Gzip)
			}
		} else {
			zap.L().Info("Handled by try router", zap.String("path", path))
			TryRootRouter(response, request, path, config.Static, config.Gzip)
		}
	}
}

func CustomRouter(response http.ResponseWriter, request *http.Request, URLPath string, staticDir string, cus []utils.CustomConfig, gzt bool) bool {
	if len(cus) == 0 {
		return false
	}

	zap.L().Info("Handled by custom router", zap.String("path", URLPath))
	for _, custom := range cus {
		if custom.Urlpath == URLPath {
			fullPath := filepath.Join(staticDir, filepath.Clean(custom.Filepath))
			SendStaticFile(response, request, fullPath, gzt)
			return true
		}
	}
	return false
}

func Router(response http.ResponseWriter, request *http.Request, URLPath string, h utils.HtmlConfig, gzt bool) {
	switch URLPath {
	case "/":
		fullPath := filepath.Join(h.Dirpath, h.Index)
		SendStaticFile(response, request, fullPath, gzt)
	default:
		fullPath := filepath.Join(h.Dirpath, filepath.Clean(URLPath))
		SendStaticFile(response, request, fullPath, gzt)
	}
	zap.L().Info("Serving static file", zap.String("path", URLPath))
}

func TryRootRouter(response http.ResponseWriter, request *http.Request, URLPath string, h utils.HtmlConfig, gzt bool) {
	switch URLPath {
	case "/":
		fullPath := filepath.Join(h.Dirpath, h.Index)
		SendTryRootFile(response, request, fullPath, h, gzt)
	default:
		fullPath := filepath.Join(h.Dirpath, filepath.Clean(URLPath))
		SendTryRootFile(response, request, fullPath, h, gzt)
	}
	zap.L().Info("Serving try static file", zap.String("path", URLPath))
}
