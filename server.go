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
	r := bufio.NewScanner(conn)
	w := bufio.NewWriter(conn)
	for {
		if r.Scan() == false {
			if r.Err() == nil {
				log.Printf("disconnect %v", conn.RemoteAddr().String())
			} else {
				log.Printf("recv error %v %v", conn.RemoteAddr().String(), r.Err())
			}
			return
		}
		go func(text string) {
			w.WriteString(text + "\n")
			w.Flush()
		}(r.Text())
	}
}
