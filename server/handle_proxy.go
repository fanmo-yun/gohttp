package server

import (
	"gohttp/utils"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

type LoadBalancer struct {
	Backends []*httputil.ReverseProxy
	Current  uint32
}

func NewLoadBalancer(backends []*httputil.ReverseProxy) *LoadBalancer {
	return &LoadBalancer{Backends: backends, Current: 0}
}

func (lb *LoadBalancer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	backend := lb.nextBackend()
	if backend != nil {
		backend.ServeHTTP(res, req)
	} else {
		http.Error(res, "Service unavailable", http.StatusServiceUnavailable)
	}
}

func (lb *LoadBalancer) nextBackend() *httputil.ReverseProxy {
	n := atomic.AddUint32(&lb.Current, 1)
	index := (int(n) - 1) % len(lb.Backends)

	if index < 0 || index >= len(lb.Backends) {
		log.Println("Invalid backend index calculated:", index)
		return nil
	}

	return lb.Backends[index]
}

func NewProxy(target string) (*httputil.ReverseProxy, error) {
	targetURL, parseErr := url.Parse(target)
	if parseErr != nil {
		return nil, parseErr
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		SendHTTPErrorResponse(w, http.StatusBadGateway)
	}
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
			continue
		}
		proxies[rp.PathPrefix] = p
	}
	return proxies
}

func FindAndServeProxy(response http.ResponseWriter, request *http.Request, path string, proxies map[string]*httputil.ReverseProxy) bool {
	for prefix, proxy := range proxies {
		if len(path) >= len(prefix) && path[:len(prefix)] == prefix {
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
			continue
		}
		urls = append(urls, backend)
	}

	if len(urls) == 0 {
		return nil
	}
	return urls
}
