package proxy

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

const (
	testMsg    = "test message"
	serverPort = ":5000"
	proxyPort  = ":3000"
	connType   = "tcp"
	pass       = "pass"
	bufSize    = 256
)

func server(t *testing.T) {
	ln, err := net.Listen(connType, serverPort)
	assert.Nil(t, err)
	defer ln.Close()

	conn, err := ln.Accept()
	assert.Nil(t, err)
	defer conn.Close()

	buf := make([]byte, bufSize)
	_, err = conn.Read(buf)
	assert.Nil(t, err)

	assert.Equal(t, testMsg, string(bytes.Trim(buf, "\x00")))
}

func client(t *testing.T) {
	conn, err := net.Dial(connType, proxyPort)
	assert.Nil(t, err)
	defer conn.Close()

	_, err = conn.Write([]byte(testMsg))
	assert.Nil(t, err)
}

func TestProxyDefault(t *testing.T) {
	go client(t)

	go func() {
		cfg := Proxy{
			From: proxyPort,
			To:   serverPort,
		}
		server := NewProxyServer(cfg)

		server.Start()
	}()

	server(t)
}
