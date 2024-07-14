package server

import (
	"gohttp/utils"
	"net/http"
	"net/url"
	"path/filepath"
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
		urlPath := Handle(response, request)
		if urlPath == nil {
			return
		}
		path := urlPath.Path

		if proxies != nil && FindAndServeProxy(response, request, path, proxies) {
			return
		}

		if backends != nil {
			loadBalancer.ServeHTTP(response, request)
			return
		}

		if !CustomRouter(response, request, path, config.Static.Dirpath, config.Custom) {
			Router(response, request, path, config.Static)
		}
	}
}

func CustomRouter(response http.ResponseWriter, request *http.Request, URLPath string, staticDir string, cus []utils.CustomConfig) bool {
	if len(cus) == 0 {
		return false
	}

	for _, custom := range cus {
		if custom.Urlpath == URLPath {
			fullPath := filepath.Join(staticDir, filepath.Clean(custom.Filepath))
			SendStaticFile(response, request, fullPath)
			return true
		}
	}
	return false
}

func Router(response http.ResponseWriter, request *http.Request, URLPath string, h utils.HtmlConfig) {
	switch URLPath {
	case "/":
		fullPath := filepath.Join(h.Dirpath, h.Index)
		SendStaticFile(response, request, fullPath)
	default:
		fullPath := filepath.Join(h.Dirpath, filepath.Clean(URLPath))
		SendStaticFile(response, request, fullPath)
	}
}
