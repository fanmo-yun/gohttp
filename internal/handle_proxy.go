package server

import (
	"gohttpd/utils"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"

	"go.uber.org/zap"
)

type LoadBalancer struct {
	Backends []*httputil.ReverseProxy
	Current  uint32
}

func NewLoadBalancer(backends []*httputil.ReverseProxy) *LoadBalancer {
	return &LoadBalancer{Backends: backends, Current: 0}
}

func (lb *LoadBalancer) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	backend := lb.nextBackend()
	if backend != nil {
		zap.L().Info("Forwarding request to backend")
		backend.ServeHTTP(response, request)
	} else {
		zap.L().Error("No available backend to handle the request")
		SendHTTPErrorResponse(response, http.StatusServiceUnavailable)
	}
}

func (lb *LoadBalancer) nextBackend() *httputil.ReverseProxy {
	n := atomic.AddUint32(&lb.Current, 1)
	index := (int(n) - 1) % len(lb.Backends)

	if index < 0 || index >= len(lb.Backends) {
		zap.L().Error("Invalid backend index calculated", zap.Int("index", index))
		return nil
	}

	return lb.Backends[index]
}

func NewProxy(target string) (*httputil.ReverseProxy, error) {
	targetURL, parseErr := url.Parse(target)
	if parseErr != nil {
		zap.L().Error("Failed to parse target URL", zap.String("target", target), zap.Error(parseErr))
		return nil, parseErr
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		zap.L().Error("Error encountered while handling request", zap.String("target", target), zap.Error(err))
		SendHTTPErrorResponse(w, http.StatusBadGateway)
	}
	zap.L().Info("Created new reverse proxy", zap.String("target", target))
	return proxy, nil
}

func CreateProxies(reverseproxies []utils.ProxyConfig) map[string]*httputil.ReverseProxy {
	if len(reverseproxies) == 0 {
		return nil
	}

	proxies := make(map[string]*httputil.ReverseProxy)
	for _, rp := range reverseproxies {
		p, err := NewProxy(rp.TargetURL)
		if err != nil {
			zap.L().Error("Failed to create proxy", zap.String("target", rp.TargetURL), zap.Error(err))
			continue
		}
		proxies[rp.PathPrefix] = p
		zap.L().Info("Created proxy for path prefix", zap.String("pathPrefix", rp.PathPrefix), zap.String("target", rp.TargetURL))
	}
	return proxies
}

func FindAndServeProxy(response http.ResponseWriter, request *http.Request, URLPath string, proxies map[string]*httputil.ReverseProxy) bool {
	for prefix, proxy := range proxies {
		if len(URLPath) >= len(prefix) && URLPath[:len(prefix)] == prefix {
			zap.L().Info("Proxying request", zap.String("path", URLPath), zap.String("prefix", prefix))
			proxy.ServeHTTP(response, request)
			return true
		}
	}
	return false
}

func NewBackend(backendUrl []utils.BackendConfig) []*httputil.ReverseProxy {
	if len(backendUrl) == 0 {
		return nil
	}

	urls := []*httputil.ReverseProxy{}
	for _, url := range backendUrl {
		backend, err := NewProxy(url.BackendURL)
		if err != nil {
			zap.L().Error("Failed to create backend proxy", zap.String("backendURL", url.BackendURL), zap.Error(err))
			continue
		}
		urls = append(urls, backend)
		zap.L().Info("Created backend proxy", zap.String("backendURL", url.BackendURL))
	}

	if len(urls) == 0 {
		zap.L().Warn("No valid backends configured")
		return nil
	}
	return urls
}
