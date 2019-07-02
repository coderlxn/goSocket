package main

import (
	"log"
	"net"
)

func send(conn net.Conn) {
	words := "Hello World!"
	conn.Write([]byte(words))
	log.Println("send finished")
}

func main()  {
	server := "127.0.0.1:2048"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		log.Println(err.Error())
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("connect success")
	send(conn)
	conn.Close()
}