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

func TestDefault(t *testing.T) {
	go func() {
		conn, err := net.Dial(connType, proxyPort)
		assert.Nil(t, err)
		defer conn.Close()

		_, err = conn.Write([]byte(testMsg))
		assert.Nil(t, err)
	}()

	var proxy *Proxy

	go func() {
		cfg := Proxy{
			From: proxyPort,
			To:   serverPort,
		}
		proxy = NewProxyServer(cfg)

		proxy.Start()
	}()

	server(t)
}

func TestCorrectPass(t *testing.T) {
	go func() {
		conn, err := net.Dial(connType, proxyPort)
		assert.Nil(t, err)
		defer conn.Close()

		_, err = conn.Write([]byte(pass))
		assert.Nil(t, err)

		buf := make([]byte, bufSize)
		_, err = conn.Read(buf)
		assert.Nil(t, err)

		assert.Equal(t, "authorized", string(bytes.Trim(buf, "\x00")))

		_, err = conn.Write([]byte(testMsg))
		assert.Nil(t, err)
	}()

	go func() {
		cfg := Proxy{
			From:     proxyPort,
			To:       serverPort,
			Password: []byte(pass),
		}
		server := NewProxyServer(cfg)

		server.Start()
	}()

	server(t)
}

func TestIncorrectPass(t *testing.T) {
	go func() {
		conn, err := net.Dial(connType, proxyPort)
		assert.Nil(t, err)
		defer conn.Close()

		_, err = conn.Write([]byte("lal"))
		assert.Nil(t, err)

		buf := make([]byte, bufSize)
		_, err = conn.Read(buf)
		assert.Nil(t, err)

		assert.Equal(t, "authorized", string(bytes.Trim(buf, "\x00")))
	}()

	cfg := Proxy{
		From:     proxyPort,
		To:       serverPort,
		Password: []byte(pass),
	}
	server := NewProxyServer(cfg)

	server.Start()
}
