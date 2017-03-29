# TCP proxy server
TCP proxy server
## Usage
### CLI
#### Start proxy server
```
NAME:
   proxy.go - TCP proxy server

USAGE:
   proxy.go [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   -f value, --from value  set proxy server port
   -t value, --to value    set proxy server redirect port
   -p value, --pass value  set proxy server password
   -l, --log               enable logging
   -b value, --buf value   set buffer size (default: 256)
   --help, -h              show help
   --version, -v           print the version
```
```
$ go run proxy.go --from :3333 --to :9999
```
### GO
```go
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
```
