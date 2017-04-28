package main

import (
	"bufio"
	"log"
	"net"
)

var clients map[*net.TCPConn]struct{}

func main() {
	addr, err := net.ResolveTCPAddr("tcp4", ":1234")
	if err != nil {
		panic(err)
	}
	listen, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		panic(err)
	}
	log.Printf("listen %v", addr)

	clients = make(map[*net.TCPConn]struct{})

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			log.Printf("accept error %v", err)
			continue
		}
		log.Printf("accept %v", conn.RemoteAddr().String())
		go handleClient(conn)
	}
}

func handleClient(conn *net.TCPConn) {
	defer conn.Close()
	clients[conn] = struct{}{}
	defer delete(clients, conn)
	r := bufio.NewScanner(conn)
	for {
		if r.Scan() == false {
			if r.Err() == nil {
				log.Printf("disconnect %v", conn.RemoteAddr().String())
			} else {
				log.Printf("recv error %v %v", conn.RemoteAddr().String(), r.Err())
			}
			return
		}
		for c := range clients {
			if c != conn {
				go func(c *net.TCPConn, text string) {
					w := bufio.NewWriter(c)
					w.WriteString(text + "\n")
					w.Flush()
				}(c, r.Text())
			}
		}
	}
}
