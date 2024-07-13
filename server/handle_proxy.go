package server

import (
	"gohttp/utils"
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
	return &LoadBalancer{Backends: backends}
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend := lb.nextBackend()
	backend.ServeHTTP(w, r)
}

func (lb *LoadBalancer) nextBackend() *httputil.ReverseProxy {
	n := atomic.AddUint32(&lb.Current, 1)
	return lb.Backends[(int(n)-1)%len(lb.Backends)]
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

func FindAndServeProxy(res http.ResponseWriter, req *http.Request, path string, proxies map[string]*httputil.ReverseProxy) bool {
	for prefix, proxy := range proxies {
		if len(path) >= len(prefix) && path[:len(prefix)] == prefix {
			proxy.ServeHTTP(res, req)
			return true
		}
	}
	return false
}

func NewBackend(backendUrl []utils.BackendConfig) []*httputil.ReverseProxy {
	if len(backendUrl) == 0 {
		return nil
	}

	urls := make([]*httputil.ReverseProxy, 0)
	for _, url := range backendUrl {
		backend, err := NewProxy(url.BackendURL)
		if err != nil {
			continue
		}
		urls = append(urls, backend)
	}
	return urls
}
