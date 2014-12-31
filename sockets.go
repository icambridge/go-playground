package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", ":2000")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func(c net.Conn) {
			buf := make([]byte, 4096)

			for {
				n, err := c.Read(buf)
				if err != nil || n == 0 {
					c.Close()
					break
				}
				b := buf[:n]
				s := string(b)
				s = strings.TrimSpace(s)
				fmt.Print(s)
				fmt.Println("------")
				n, err = c.Write(b)
				if err != nil {
					c.Close()
					break
				}
			}
		}(conn)
	}
}
