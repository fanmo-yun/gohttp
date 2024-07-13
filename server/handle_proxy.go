package server

import (
	"gohttp/utils"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewProxy(target string) (*httputil.ReverseProxy, error) {
	targetURL, parseErr := url.Parse(target)
	if parseErr != nil {
		return nil, parseErr
	}
	return httputil.NewSingleHostReverseProxy(targetURL), nil
}

func CreateProxies(reverseproxies []utils.ProxyConfig) map[string]*httputil.ReverseProxy {
	if len(reverseproxies) == 0 {
		return nil
	}

	proxies := make(map[string]*httputil.ReverseProxy)
	for _, rp := range reverseproxies {
		p, err := NewProxy(rp.TargetURL)
		if err != nil {
			continue
		}
		proxies[rp.PathPrefix] = p
	}
	return proxies
}

func FindAndServeProxy(res http.ResponseWriter, req *http.Request, path string, proxies map[string]*httputil.ReverseProxy) bool {
	for prefix, proxy := range proxies {
		if len(path) >= len(prefix) && path[:len(prefix)] == prefix {
			proxy.ServeHTTP(res, req)
			return true
		}
	}
	return false
}
