package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

type Backend struct {
	URL          *url.URL
	ReverseProxy *httputil.ReverseProxy
}

type LoadBalancer struct {
	backends []*Backend
	current  uint64
}

func (lb *LoadBalancer) getNextBackend() *Backend {
	index := atomic.AddUint64(&lb.current, 1)
	return lb.backends[int(index)%len(lb.backends)]
}

func (lb *LoadBalancer) handler(w http.ResponseWriter, r *http.Request) {
	backend := lb.getNextBackend()
	fmt.Println("Routing to :- %s", backend.URL.String())
	backend.ReverseProxy.ServeHTTP(w, r)
}

func main() {
	backendUrls := []string{
		"http://localhost:8081",
		"http://localhost:8082",
	}

	var backends []*Backend
	for _, addr := range backendUrls {
		parsedURL, err := url.Parse(addr)
		if err != nil {
			fmt.Printf("invalid url %v: %v", addr, err)
		}

		proxy := httputil.NewSingleHostReverseProxy(parsedURL)
		backends = append(backends, &Backend{
			URL:          parsedURL,
			ReverseProxy: proxy,
		})
	}

	lb := &LoadBalancer{backends: backends}

	fmt.Println("Load Balancer running on 8083")
	if err := http.ListenAndServe(":8083", http.HandlerFunc(lb.handler)); err != nil {
		fmt.Printf("load balancer error :- %v", err)
	}
}
