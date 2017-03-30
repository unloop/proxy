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
	ln       net.Listener
	target   net.Conn
	conns    []net.Conn
	quit     chan bool
}

func NewProxyServer(cfg Proxy) *Proxy {
	var bs []byte

	r, err := regexp.Compile(":[\\d]{4}")
	check(err)
	if !r.MatchString(cfg.From) || (!r.MatchString(cfg.To)) {
		log.Fatal("incorrect ports")
	}

	if len(cfg.Password) != 0 {
		bs = encode(cfg.Password)
	} else {
		bs = cfg.Password
	}

	if cfg.BufSize == 0 {
		cfg.BufSize = 256
	}

	return &Proxy{
		From:     cfg.From,
		To:       cfg.To,
		Logging:  cfg.Logging,
		Password: bs,
		BufSize:  cfg.BufSize,
		quit:     make(chan bool, 1),
	}
}

func (p *Proxy) Start() {
	if p.started {
		log.Panic("proxy server already started")
	}

	p.started = true

	var err error

	p.ln, err = net.Listen("tcp", p.From)
	if err != nil {
		log.Panic(err)
	}
	defer p.ln.Close()

	p.pLog("starting proxy on port " + p.From[1:] + " forwarding to " + p.To[1:])

	for {
		if conn, err := p.ln.Accept(); err == nil {
			go p.nClient(conn)

			p.conns = append(p.conns, conn)
		} else {
			select {
			case <-p.quit:
				return
			default:
			}

			log.Panic("error accept listen")
		}
	}
}

func (p *Proxy) nClient(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, p.BufSize)

	if len(p.Password) != 0 {
		if n, err := conn.Read(buf); err == nil && n > 0 {
			bs := encode(bytes.Trim(buf, "\x00"))

			if !bytes.Equal(bs, p.Password) {
				p.incorrectPass(conn)
				return
			}

			_, err = conn.Write([]byte("authorized"))
			check(err)

			p.pLog(conn.RemoteAddr().String() + " authorized: OK")
		} else {
			p.incorrectPass(conn)
			return
		}
	} else {
		p.pLog("new client " + conn.RemoteAddr().String())
	}

	var err error

	p.target, err = net.Dial("tcp", p.To)
	if err != nil {
		p.pLog("error to start dial connection to " + p.To)
		conn.Close()
		return
	}
	defer p.target.Close()

	for {
		buf = make([]byte, p.BufSize)

		n, err := conn.Read(buf)
		if err != nil {
			_, ok := err.(net.Error)
			if err == io.EOF || ok {
				p.pLog("client " + conn.RemoteAddr().String() + " close connection")
				return
			} else {
				check(err)
			}
		}

		if n > 0 {
			p.pLog(
				"message \"" + string(bytes.Trim(buf, "\x00")) +
					"\" received from " + conn.RemoteAddr().String() +
					" forwarded to " + p.target.RemoteAddr().String(),
			)

			_, err = p.target.Write(buf)
			check(err)
		}
	}
}

func (p *Proxy) Close() {
	close(p.quit)

	p.ln.Close()
	for _, conn := range p.conns {
		conn.Close()
	}

	p.pLog("close proxy with all connections")
}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func (p *Proxy) incorrectPass(conn net.Conn) {
	_, err := conn.Write([]byte("Incorrect password, connection close"))
	check(err)

	p.pLog(conn.RemoteAddr().String() + " incorrect pass")
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
