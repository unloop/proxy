package proxy

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
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
	defer ln.Close()

	conn, err := ln.Accept()
	defer conn.Close()

	buf := make([]byte, bufSize)
	_, err = conn.Read(buf)
	assert.NoError(t, err)

	assert.Equal(t, testMsg, string(bytes.Trim(buf, "\x00")))
}

func TestDefault(t *testing.T) {
	go func() {
		conn, err := net.Dial(connType, proxyPort)
		assert.NoError(t, err)
		defer conn.Close()

		_, err = conn.Write([]byte(testMsg))
		assert.NoError(t, err)
	}()

	var (
		proxy *Proxy
		err   error
	)

	go func() {
		proxy, err = NewProxyServer(Proxy{
			From: proxyPort,
			To:   serverPort,
		})
		assert.NoError(t, err)

		err = proxy.Start()
		assert.NoError(t, err)
	}()

	go server(t)

	time.Sleep(time.Second)
	proxy.Close()
}

func TestCorrectPass(t *testing.T) {
	go func() {
		conn, err := net.Dial(connType, proxyPort)
		assert.NoError(t, err)
		defer conn.Close()

		_, err = conn.Write([]byte(pass))
		assert.NoError(t, err)

		buf := make([]byte, bufSize)
		_, err = conn.Read(buf)
		assert.NoError(t, err)

		assert.Equal(t, "authorized", string(bytes.Trim(buf, "\x00")))

		_, err = conn.Write([]byte(testMsg))
		assert.NoError(t, err)
	}()

	var (
		proxy *Proxy
		err   error
	)

	go func() {
		proxy, err = NewProxyServer(Proxy{
			From:     proxyPort,
			To:       serverPort,
			Password: []byte(pass),
		})
		assert.NoError(t, err)

		err = proxy.Start()
		assert.NoError(t, err)
	}()

	go server(t)

	time.Sleep(time.Second)
	proxy.Close()
}

func TestIncorrectPass(t *testing.T) {
	go func() {
		conn, err := net.Dial(connType, proxyPort)
		assert.NoError(t, err)
		defer conn.Close()

		_, err = conn.Write([]byte("incorrect pass"))
		assert.NoError(t, err)

		buf := make([]byte, bufSize)
		_, err = conn.Read(buf)
		assert.NoError(t, err)

		assert.Equal(t, "incorrect password, connection close", string(bytes.Trim(buf, "\x00")))
	}()

	var (
		proxy *Proxy
		err   error
	)

	go func() {
		proxy, err = NewProxyServer(Proxy{
			From:     proxyPort,
			To:       serverPort,
			Password: []byte(pass),
		})
		assert.NoError(t, err)

		err = proxy.Start()
		assert.NoError(t, err)
	}()

	time.Sleep(time.Second)
	proxy.Close()
}

func TestAlreadyStarted(t *testing.T) {
	var (
		proxy *Proxy
		err   error
	)
	pending := make(chan bool, 1)

	go func() {
		proxy, err = NewProxyServer(Proxy{
			From: proxyPort,
			To:   serverPort,
		})
		assert.NoError(t, err)

		close(pending)

		err = proxy.Start()
		assert.NoError(t, err)
	}()

	<-pending
	err = proxy.Start()
	assert.EqualError(t, err, "proxy server already started")

	time.Sleep(time.Second)
	proxy.Close()
}

func TestIncorrectPorts(t *testing.T) {
	_, err := NewProxyServer(Proxy{
		From: "incorrect ports",
		To:   "incorrect ports",
	})

	assert.EqualError(t, err, "entered incorrect ports")
}
