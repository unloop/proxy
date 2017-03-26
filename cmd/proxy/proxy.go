package main

import (
	"github.com/lavrs/proxy"
)

func main() {
	cfg := proxy.Proxy{
		From:     ":3333",
		To:       ":9999",
		Logging:  true,
		Password: []byte("password"),
	}
	server := proxy.NewProxyServer(cfg)

	server.Start()
}
