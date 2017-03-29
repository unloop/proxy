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
	ch := make(chan os.Signal)

	go func() {
		conn, err = net.Dial("tcp", ":3333")
		if err != nil {
			log.Panic(err)
		}

		for {
			fmt.Print("text: ")
			fmt.Scanln(&buf)

			_, err = conn.Write([]byte(buf))
			if err != nil {
				log.Panic(err)
			}
		}
	}()

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	<-ch
}
