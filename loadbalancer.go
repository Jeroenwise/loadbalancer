package main

import (
	"net/http"
	"net/http/httputil"
)

var servers = []string{
	"localhost:3000",
	"localhost:3001",
	"localhost:3002",
}

type serverProvider func() string

func roundRobinProvider(hosts []string) serverProvider {
	var balanceIndex = 0

	return func() string {
		server := hosts[balanceIndex]
		balanceIndex++
		if balanceIndex >= len(hosts) {
			balanceIndex = 0
		}
		return server
	}
}

func newLoadBalancedReverseProxy(provider serverProvider) *httputil.ReverseProxy {
	director := func(request *http.Request) {
		request.URL.Scheme = "http"
		request.URL.Host = provider()
	}
	return &httputil.ReverseProxy{Director: director}
}

func main() {
	serverProvider := roundRobinProvider(servers)

	reverseProxy := newLoadBalancedReverseProxy(serverProvider)

	http.ListenAndServe(":8080", reverseProxy)
}
