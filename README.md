# Simple TCP proxy server
Simple TCP proxy server with logging and password

## Usage
### CLI
#### Start proxy server
```
$ go run proxy.go --from :3333 --to :9999
```
#### Start proxy server with logging and password
```
$ go run proxy.go -f :3333 -t :9999 -l --pass password
```
### GO
#### Configure proxy server
```go
cfg := proxy.Proxy{
        // proxy server port
	From:     ":3333",
	
	// forwarding port
	To:       ":9999",
	
	// enable logging
	Logging:  true,
	
	// set password
	Password: []byte("password")),
}
```
#### Start proxy server
 ```go
server := proxy.NewProxyServer(cfg)
server.Start()
```
#### Example usage
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