# TCP proxy
Simple TCP proxy server written in Go (Golang) with logging and auth
## Usage
### CLI
```
NAME:
   proxy.exe - TCP proxy server

USAGE:
   proxy.exe [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   -f value, --from value  set proxy server port
   -t value, --to value    set proxy server redirect port
   -a value, --auth value  set proxy server auth token
   -l, --log               enable logging
   --help, -h              show help
   --version, -v           print the version
```
##### Start proxy server
```
$ go run proxy.go -l -f :3000 -t :5000
```
### GO
```go
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
```
