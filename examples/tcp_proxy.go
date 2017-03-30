package main

import (
	"github.com/lavrs/proxy"
	"log"
)

func main() {
	server, err := proxy.NewProxyServer(proxy.Proxy{
		From:    ":3000",
		To:      ":5000",
		Logging: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
