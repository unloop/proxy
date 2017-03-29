package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var conn net.Conn
	var err error

	ch := make(chan os.Signal)

	go func() {
		conn, err = net.Dial("tcp", ":3333")
		if err != nil {
			log.Panic(err)
		}

		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("text to send: ")
			text, _, err := reader.ReadLine()
			if err != nil {
				log.Panic(err)
			}

			conn.Write([]byte(text))
		}
	}()

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	conn.Close()
}
