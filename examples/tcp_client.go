package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var (
		conn net.Conn
		err  error
		buf  string
	)
	pending := make(chan os.Signal)

	go func() {
		conn, err = net.Dial("tcp", ":3000")
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		for {
			fmt.Print("send message: ")
			fmt.Scanln(&buf)

			_, err = conn.Write([]byte(buf))
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	signal.Notify(pending, syscall.SIGINT, syscall.SIGTERM)
	<-pending

	conn.Close()
}
