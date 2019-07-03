package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

const HeaderText = "Header"
const HeaderTextLength = len(HeaderText)
const LengthTextLength = 5

func main() {
	netListen, err := net.Listen("tcp", "localhost:2048")
	if err != nil {
		log.Println(err)
		return
	}
	defer netListen.Close()

	log.Println("Waiting for clients")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println(conn.RemoteAddr().String(), " tcp connect success")
		handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	buffer := make([]byte, 2048, 2048)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println(conn.RemoteAddr().String(), " read error: ", err)
			return
		}
		tmpBuffer := buffer[:n]
		log.Println(conn.RemoteAddr().String(), "receive data string:")
		//解析buffer中的内容
		for {
			if len(tmpBuffer) == 0 {
				break
			}

			if string(tmpBuffer[:HeaderTextLength]) != HeaderText {
				panic("buffer not started with 'Header'")
			}
			lengthText := string(tmpBuffer[HeaderTextLength: HeaderTextLength + LengthTextLength])
			textLength, err := strconv.Atoi(lengthText)
			if err != nil {
				log.Println(err)
			}
			content := tmpBuffer[HeaderTextLength + LengthTextLength: HeaderTextLength + LengthTextLength + textLength]
			tmpBuffer = tmpBuffer[HeaderTextLength + LengthTextLength + textLength:]

			fmt.Println(string(content))
		}

	}
}
