package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("launching server..")

	var (
		conns []net.Conn
		ln    net.Listener
	)
	pending := make(chan os.Signal)

	go func() {
		ln, err := net.Listen("tcp", ":5000")
		if err != nil {
			log.Fatal(err)
		}
		defer ln.Close()

		for {
			if conn, err := ln.Accept(); err == nil {
				conns = append(conns, conn)

				go func(conn net.Conn) {
					defer conn.Close()

					buf := make([]byte, 10240)

					for {
						buf = make([]byte, 10240)

						n, err := conn.Read(buf)
						if err != nil {
							if err == io.EOF {
								log.Println("lost connection to:", conn.RemoteAddr())
								return
							} else {
								log.Fatal(err)
							}
						} else if n > 0 {
							log.Println("new message: \""+string(bytes.Trim(buf, "\x00"))+
								"\" from:", conn.RemoteAddr().String())
						}
					}
				}(conn)
			} else {
				log.Fatal(err)
			}
		}
	}()

	signal.Notify(pending, syscall.SIGINT, syscall.SIGTERM)
	<-pending

	for _, conn := range conns {
		conn.Close()
	}

	ln.Close()
}
