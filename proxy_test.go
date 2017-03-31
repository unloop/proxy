package proxy

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

const (
	message = "message"

	serverPort = ":5000"
	proxyPort  = ":3000"
	connType   = "tcp"

	auth    = "token"
	notAuth = ""

	bufSize = 256

	testForwarding = "testForwarding"
	testAuthOK     = "testAuthOK"
	testAuthFail   = "testAuthFail"
)

func client(t *testing.T, cType string) {
	conn, err := net.Dial(connType, proxyPort)
	assert.NoError(t, err)
	defer conn.Close()

	switch cType {
	case testForwarding:
		_, err = conn.Write([]byte(message))
		assert.NoError(t, err)

		break
	case testAuthOK:
		_, err = conn.Write([]byte(auth))
		assert.NoError(t, err)

		buf := make([]byte, bufSize)
		_, err = conn.Read(buf)
		assert.NoError(t, err)

		assert.Equal(t, "authorized", string(bytes.Trim(buf, "\x00")))

		_, err = conn.Write([]byte(message))
		assert.NoError(t, err)

		break
	case testAuthFail:
		_, err = conn.Write([]byte("bad token"))
		assert.NoError(t, err)

		buf := make([]byte, bufSize)
		_, err = conn.Read(buf)
		assert.NoError(t, err)

		assert.Equal(t, "auth fail: connection close", string(bytes.Trim(buf, "\x00")))

		break
	default:
		t.Fatal("incorrect cType")
	}
}

func server(t *testing.T, sType bool, sReady, done chan bool) {
	ln, err := net.Listen(connType, serverPort)
	assert.NoError(t, err)
	defer func() {
		ln.Close()

		done <- true
	}()

	sReady <- true

	if sType {
		conn, err := ln.Accept()
		assert.NoError(t, err)
		defer conn.Close()

		buf := make([]byte, bufSize)
		_, err = conn.Read(buf)
		assert.NoError(t, err)

		assert.Equal(t, message, string(bytes.Trim(buf, "\x00")))
	}
}

func proxy(t *testing.T, auth string, pDone chan bool, pReady chan *Proxy) {
	proxy, err := NewProxyServer(Proxy{
		From:  proxyPort,
		To:    serverPort,
		Token: []byte(auth),
	})
	assert.NoError(t, err)

	pReady <- proxy

	err = proxy.Start()
	assert.NoError(t, err)

	pDone <- true
}

func TestForwarding(t *testing.T) {
	pReady := make(chan *Proxy)
	sReady := make(chan bool)
	done := make(chan bool)
	pDone := make(chan bool)

	go proxy(t, notAuth, pDone, pReady)
	pObj := <-pReady

	go server(t, true, sReady, done)
	<-sReady

	go client(t, testForwarding)

	<-done

	pObj.Close()

	<-pDone

	close(done)
	close(pReady)
	close(sReady)
	close(pDone)
}

func TestAuthOK(t *testing.T) {
	pReady := make(chan *Proxy)
	sReady := make(chan bool)
	done := make(chan bool)
	pDone := make(chan bool)

	go proxy(t, auth, pDone, pReady)
	pObj := <-pReady

	go server(t, true, sReady, done)
	<-sReady

	go client(t, testAuthOK)

	<-done

	pObj.Close()

	<-pDone

	close(done)
	close(pReady)
	close(sReady)
	close(pDone)
}

func TestAuthFail(t *testing.T) {
	pReady := make(chan *Proxy)
	sReady := make(chan bool)
	done := make(chan bool)
	pDone := make(chan bool)

	go proxy(t, auth, pDone, pReady)
	pObj := <-pReady

	go server(t, false, sReady, done)
	<-sReady

	client(t, testAuthFail)

	<-done

	pObj.Close()

	<-pDone

	close(done)
	close(pReady)
	close(sReady)
	close(pDone)
}

func TestAlreadyStarted(t *testing.T) {
	pReady := make(chan *Proxy)
	pDone := make(chan bool)

	go proxy(t, notAuth, pDone, pReady)
	pObj := <-pReady

	time.Sleep(time.Millisecond)

	err := pObj.Start()
	assert.EqualError(t, err, "proxy server already started")

	pObj.Close()

	<-pDone

	close(pReady)
	close(pDone)
}

func TestIncorrectPorts(t *testing.T) {
	_, err := NewProxyServer(Proxy{
		From: "bad port",
		To:   "bad port",
	})

	assert.EqualError(t, err, "entered incorrect ports")
}
