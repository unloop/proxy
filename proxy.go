package proxy

import (
	"bytes"
	"crypto/sha1"
	"io"
	"log"
	"net"
	"regexp"
)

type Proxy struct {
	From, To string
	Logging  bool
	Password []byte
	BufSize  int64
	started  bool
}

func NewProxyServer(cfg Proxy) *Proxy {
	var bs []byte

	r, err := regexp.Compile(":[\\d]{4}")
	check(err)
	if !r.MatchString(cfg.From) || (!r.MatchString(cfg.To)) {
		log.Panic("incorrect ports")
	}

	if len(cfg.Password) != 0 {
		bs = encode(cfg.Password)
	} else {
		bs = cfg.Password
	}

	if cfg.BufSize == 0 {
		cfg.BufSize = 10240
	}

	return &Proxy{
		From:     cfg.From,
		To:       cfg.To,
		Logging:  cfg.Logging,
		Password: bs,
		BufSize:  cfg.BufSize,
	}
}

func (p *Proxy) Start() {
	if p.started {
		log.Panic("proxy server already started")
	}

	p.started = true

	if listen, err := net.Listen("tcp", p.From); err == nil {
		defer listen.Close()

		p.pLog("starting the server on port " + p.From[1:] + " forwarding to " + p.To[1:])

		for {
			if conn, err := listen.Accept(); err == nil {
				go p.nClient(conn)
			} else {
				log.Panic("error listen.Accept")
			}
		}
	} else {
		log.Panic("error to start listen")
	}
}

func (p *Proxy) nClient(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, p.BufSize)

	if len(p.Password) != 0 {
		if n, err := conn.Read(buf); (err == nil) && (n > 0) {
			bs := encode(bytes.Trim(buf, "\x00"))

			if !bytes.Equal(bs, p.Password) {
				p.incorrectPass(conn)
				return
			}
		} else {
			p.incorrectPass(conn)
			return
		}
	}

	p.pLog("new client " + conn.RemoteAddr().String())

	if target, err := net.Dial("tcp", p.To); err == nil {
		defer target.Close()

		for {
			buf = make([]byte, p.BufSize)

			n, err := conn.Read(buf)
			if err != nil {
				if err == io.EOF {
					p.pLog("client " + conn.RemoteAddr().String() + " close connection")
					return
				} else {
					check(err)
				}
			}

			buf = bytes.Trim(buf, "\x00")

			if n > 0 {
				p.pLog(
					"message " + string(bytes.Trim(buf, "\x00")) +
						" received from " + conn.RemoteAddr().String() +
						" forwarded to " + target.RemoteAddr().String(),
				)

				_, err = target.Write(buf)
				check(err)
			}
		}
	} else {
		log.Panic("error to start dial connection")
	}
}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func (p *Proxy) incorrectPass(conn net.Conn) {
	_, err := conn.Write([]byte("Incorrect password, connection close"))
	check(err)

	p.pLog(conn.RemoteAddr().String() + "send incorrect address")
}

func (p *Proxy) pLog(l string) {
	if p.Logging {
		log.Println(l)
	}
}

func encode(pass []byte) []byte {
	h := sha1.New()
	h.Write(pass)
	bs := h.Sum(nil)

	return bs
}
