package main

import (
	"log"
	"net"
)

func main() {
	netListen, err := net.Listen("tcp", "localhost:2048")
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer netListen.Close()

	log.Println("Waiting for clients")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		log.Println(conn.RemoteAddr().String(), " tcp connect success")
		handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	buffer := make([]byte, 2048)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println(conn.RemoteAddr().String(), " read error: ", err.Error())
			return
		}
		log.Println(conn.RemoteAddr().String(), "receive data string:\n", string(buffer[:n]))
	}
}
