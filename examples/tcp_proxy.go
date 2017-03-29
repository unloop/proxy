package main

import "github.com/lavrs/proxy"

func main() {
	cfg := proxy.Proxy{
		From:    ":3333",
		To:      ":9999",
		Logging: true,
	}
	server := proxy.NewProxyServer(cfg)

	server.Start()
}
