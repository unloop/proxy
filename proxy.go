package proxy

import (
	"log"
	"net"
	"bytes"
)

type Proxy struct {
	From, To string
	Logging  bool
}

func NewProxyServer(cfg Proxy) *Proxy {
	return &Proxy{
		From:    cfg.From,
		To:      cfg.To,
		Logging: cfg.Logging,
	}
}

func (p *Proxy) Start() error {
	log.Println("Starting the server...")

	listen, err := net.Listen("tcp", p.From)
	check(err)

	for {
		conn, err := listen.Accept()
		check(err)

		go p.newClient(conn)
	}
}

func (p *Proxy) newClient(conn net.Conn) {
	if p.Logging {
		log.Println("New client", conn.RemoteAddr())
	}

	target, err := net.Dial("tcp", p.To)
	check(err)

	buf := make([]byte, 10240)

	for {
		n, err := conn.Read(buf)
		check(err)
		buf = bytes.Trim(buf, "\x00")

		if n > 0 {
			if p.Logging {
				log.Print("Message ", string(buf),
					" received from ", conn.RemoteAddr(),
					" forwarded to ", target.RemoteAddr(),
				)
			}

			target.Write(buf)
		}
	}
}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}
