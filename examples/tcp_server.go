package main

import (
	"bytes"
	"io"
	"log"
	"net"
)

func main() {
	log.Println("launching server..")

	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Panic(err)
	}

	conn, err := ln.Accept()
	if err != nil {
		log.Panic(err)
	}

	buf := make([]byte, 10240)

	for {
		buf = make([]byte, 10240)

		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Fatal("EOF")
			} else {
				log.Fatal(err)
			}
		} else if n > 0 {
			log.Println("message " + string(bytes.Trim(buf, "\x00")))
		}
	}
}
