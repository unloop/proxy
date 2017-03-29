# Example
### Usage
Start server
```console
go run tcp_server.go
```
Start proxy
```
go run tcp_proxy.go -l -f :3333 -t :9999
```
Start client
```
go run tcp_client.go
```