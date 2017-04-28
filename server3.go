package main

import (
	"bufio"
	"log"
	"net"
)

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

	clients := make(map[*net.TCPConn]struct{})
	bcast := make(chan string)
	go handleBroadcast(clients, bcast)

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			log.Printf("accept error %v", err)
			continue
		}
		log.Printf("accept %v", conn.RemoteAddr().String())
		go func() {
			defer conn.Close()
			defer delete(clients, conn)
			clients[conn] = struct{}{}
			handleClient(conn, bcast)
		}()
	}
}

func handleClient(conn *net.TCPConn, bcast chan string) {
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
		bcast <- r.Text()
	}
}

func handleBroadcast(clients map[*net.TCPConn]struct{}, bcast chan string) {
	for {
		msg := <-bcast
		for conn := range clients {
			go func(conn *net.TCPConn, msg string) {
				w := bufio.NewWriter(conn)
				w.WriteString(msg + "\n")
				w.Flush()
			}(conn, msg)
		}
	}
}
