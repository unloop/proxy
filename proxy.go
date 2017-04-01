package proxy

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"io"
	"log"
	"net"
	"regexp"
)

type Proxy struct {
	From, To string
	Logging  bool
	Token    []byte
	started  bool
	ln       net.Listener
	target   net.Conn
	conns    []net.Conn
	close    chan bool
}

func NewProxyServer(cfg Proxy) (*Proxy, error) {
	var bs []byte

	r, err := regexp.Compile(":[\\d]{4}")
	check(err)
	if !r.MatchString(cfg.From) || (!r.MatchString(cfg.To)) {
		return nil, errors.New("entered incorrect port")
	}

	if len(cfg.Token) != 0 {
		bs = encode(cfg.Token)
	} else {
		bs = cfg.Token
	}

	return &Proxy{
		From:    cfg.From,
		To:      cfg.To,
		Logging: cfg.Logging,
		Token:   bs,
		close:   make(chan bool, 1),
		started: false,
	}, nil
}

func (p *Proxy) Start() error {
	if p.started {
		return errors.New("proxy server already started")
	}

	p.started = true

	var err error

	p.ln, err = net.Listen("tcp", p.From)
	if err != nil {
		return errors.New("error listening: " + err.Error())
	}
	defer p.ln.Close()

	p.pLog("starting proxy on port " + p.From[1:] + " forwarding to " + p.To[1:] + "\n")

	for {
		if conn, err := p.ln.Accept(); err == nil {
			go p.nClient(conn)

			p.conns = append(p.conns, conn)
		} else {
			select {
			case <-p.close:
				return nil
			default:
				return errors.New("error accepting: " + err.Error())
			}
		}
	}
}

func (p *Proxy) nClient(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 256)

	if len(p.Token) != 0 {
		if n, err := conn.Read(buf); err == nil && n > 0 {
			bs := encode(bytes.Trim(buf, "\x00"))

			if !bytes.Equal(bs, p.Token) {
				p.authFail(conn)
				return
			}

			_, err = conn.Write([]byte("authorized"))
			check(err)

			p.pLog(conn.RemoteAddr().String() + ": authorized status: OK")
		} else {
			p.authFail(conn)
			return
		}
	}

	var err error

	p.target, err = net.Dial("tcp", p.To)
	check(err)
	defer p.target.Close()

	p.pLog("new client: " + conn.RemoteAddr().String() + " connecting to: " + p.target.RemoteAddr().String())

	for {
		_, err = io.Copy(p.target, conn)

		_, ok := err.(net.Error)
		if !ok {
			check(err)
		}

		p.pLog("client: " + conn.RemoteAddr().String() +
			" lose connection with: " + p.target.RemoteAddr().String())

		return
	}
}

func (p *Proxy) Close() {
	p.started = false

	close(p.close)

	if p.ln != nil {
		p.ln.Close()
	}

	if p.target != nil {
		p.target.Close()
	}

	for _, conn := range p.conns {
		conn.Close()
	}

	p.pLog("сlosing all connections")
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (p *Proxy) authFail(conn net.Conn) {
	_, err := conn.Write([]byte("auth fail: connection close"))
	check(err)

	p.pLog(conn.RemoteAddr().String() + " authorized status: FAIL")
}

func (p *Proxy) pLog(pLog string) {
	if p.Logging {
		log.Println(pLog)
	}
}

func encode(token []byte) []byte {
	h := sha1.New()
	h.Write(token)
	bs := h.Sum(nil)

	return bs
}
