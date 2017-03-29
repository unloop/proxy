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
	BufSize  int64
}

func NewProxyServer(cfg Proxy) *Proxy {
	var bs []byte

	if len(cfg.Password) != 0 {
		bs = encode(cfg.Password)
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
	listen, err := net.Listen("tcp", p.From)
	check(err)

	if p.Logging {
		log.Println("Starting the server on port", p.From[1:], "forwarding to", p.To[1:])
	}

	for {
		conn, err := listen.Accept()
		check(err)

		go p.newClient(conn)
	}
}

func (p *Proxy) newClient(conn net.Conn) {
	buf := make([]byte, p.BufSize)

	if len(p.Password) != 0 {
		n, err := conn.Read(buf)
		check(err)

		if n > 0 {
			bs := encode(bytes.Trim(buf, "\x00"))

			if !bytes.Equal(bs, p.Password) {
				incorrectPass(conn)

				return
			}
		} else {
			incorrectPass(conn)

			return
		}
	}

	if p.Logging {
		log.Println("New client", conn.RemoteAddr())
	}

	target, err := net.Dial("tcp", p.To)
	check(err)

	for {
		buf = make([]byte, p.BufSize)

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

func incorrectPass(conn net.Conn) {
	_, err := conn.Write([]byte("Incorrect password, connection close"))
	check(err)

	conn.Close()
}

func encode(pass []byte) []byte {
	h := sha1.New()
	h.Write(pass)
	bs := h.Sum(nil)

	return bs
}
