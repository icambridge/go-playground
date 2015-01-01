package main

import (
	"fmt"
	"log"
	"net"
//	"strings"
	"time"
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
			ch := make(chan []byte)
			eCh := make(chan error)

			go func(c net.Conn, ch chan []byte, eCh chan error) {
				for {
					buf := make([]byte, 4096)
					n, err := c.Read(buf)
					if err != nil {
						eCh <- err
					}

					b := buf[:n]
					ch <- b

					if err != nil {
						eCh <- err
					}
				}

			}(c, ch, eCh)

			for {
				timer := time.NewTimer(time.Second * 5)
				t := time.Now()

				select {
				case data := <-ch:
					c.Write(data)

					fmt.Printf("Msg at %s\n", t.Local())
					timer.Reset(time.Second * 5)
				case err := <-eCh:
					fmt.Println(fmt.Sprintf("%v", err))
					c.Close()
					break
					// This will timeout on the read.
				case <- timer.C:

					fmt.Printf("Time out at %s\n", t)
					c.Write([]byte("Connect time out\n"))
					c.Close()
					return
				}
			}
		}(conn)
	}
}
