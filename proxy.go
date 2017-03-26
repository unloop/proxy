package proxy

import (
	"bytes"
	"crypto/sha1"
	"log"
	"net"
)

type Proxy struct {
	From, To string
	Logging  bool
	Password []byte
}

func NewProxyServer(cfg Proxy) *Proxy {
	var bs []byte

	if len(cfg.Password) != 0 {
		h := sha1.New()
		h.Write([]byte(cfg.Password))
		bs = h.Sum(nil)
	} else {
		bs = cfg.Password
	}

	return &Proxy{
		From:     cfg.From,
		To:       cfg.To,
		Logging:  cfg.Logging,
		Password: bs,
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
	buf := make([]byte, 10240)

	if len(p.Password) != 0 {
		n, err := conn.Read(buf)
		check(err)

		if n > 0 {
			h := sha1.New()
			h.Write(bytes.Trim(buf, "\x00"))
			bs := h.Sum(nil)

			if !bytes.Equal(bs, p.Password) {
				conn.Close()
				return
			}
		} else {
			conn.Close()
			return
		}
	}

	if p.Logging {
		log.Println("New client", conn.RemoteAddr())
	}

	target, err := net.Dial("tcp", p.To)
	check(err)

	for {
		buf = make([]byte, 10240)

		n, err := conn.Read(buf)
		check(err)

		buf = bytes.Trim(buf, "\x00")

		if n > 0 {
			if p.Logging {
				log.Print(
					"Message ", string(bytes.Trim(buf, "\x00")),
					" received from ", conn.RemoteAddr(),
					" forwarded to ", target.RemoteAddr(),
				)
			}

			_, err = target.Write(buf)
			check(err)
		}
	}
}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}
