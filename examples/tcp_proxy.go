package main

import "github.com/lavrs/proxy"

func main() {
	server := proxy.NewProxyServer(proxy.Proxy{
		From:    ":3000",
		To:      ":5000",
		Logging: true,
	})

	server.Start()
}
