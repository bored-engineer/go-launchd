package main

import (
	"io"
	"log"
	"net"

	launchd "github.com/bored-engineer/go-launchd"
)

func main() {
	l, err := launchd.Activate("Listeners")
	if err != nil {
		log.Fatalf("launchd.Socket failed: %s", err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("(net.Listener).Accept failed: %s", err)
			continue
		}
		go func(conn net.Conn) {
			defer func() {
				conn.Close()
			}()
			io.Copy(conn, conn)
		}(conn)
	}
}
