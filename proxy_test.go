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
		proxy = NewProxyServer(Proxy{
			From:    proxyPort,
			To:      serverPort,
			Logging: true,
		})
		proxy.Start()
	}()

	go server(t)

	func(){
		time.Sleep(time.Second)
		proxy.Close()
	}()
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

	var proxy *Proxy

	go func() {
		proxy = NewProxyServer(Proxy{
			From:     proxyPort,
			To:       serverPort,
			Logging:  true,
			Password: []byte(pass),
		})
		proxy.Start()
	}()

	go server(t)

	func(){
		time.Sleep(time.Second)
		proxy.Close()
	}()
}

func TestIncorrectPass(t *testing.T) {
	go func() {
		conn, err := net.Dial(connType, proxyPort)
		assert.Nil(t, err)
		defer conn.Close()

		_, err = conn.Write([]byte("incorrect pass"))
		assert.Nil(t, err)

		buf := make([]byte, bufSize)
		_, err = conn.Read(buf)
		assert.Nil(t, err)

		assert.Equal(t, "Incorrect password, connection close", string(bytes.Trim(buf, "\x00")))
	}()

	var proxy *Proxy

	go func() {
		proxy = NewProxyServer(Proxy{
			From:     proxyPort,
			To:       serverPort,
			Logging:  true,
			Password: []byte(pass),
		})
		proxy.Start()
	}()

	func(){
		time.Sleep(time.Second)
		proxy.Close()
	}()
}
