package main

import (
	"github.com/lavrs/proxy"
)

func main() {
	server := proxy.NewProxy(":3333", ":9999")

	server.Start()
}
