package proxy

import (
	"log"
	"net"
)

type proxy struct {
	From, To string
}

func NewProxy(from, to string) proxy {
	return proxy{
		From: from,
		To:   to,
	}
}

func (p *proxy) Start() error {
	log.Println("Starting the server...")

	listen, err := net.Listen("tcp", p.From)
	check(err)

	for {
		conn, err := listen.Accept()
		check(err)

		go newClient(conn, p.To)
	}
}

func newClient(listen net.Conn, to string) {
	target, err := net.Dial("tcp", to)
	check(err)

	var buf []byte

	for {
		n, err := listen.Read(buf)
		check(err)

		if n > 0 {
			target.Write(buf)
		}
	}
}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}
