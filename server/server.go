package main

import (
	"fmt"
	"github.com/labstack/gommon/random"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

const HeaderText = "Header"
const HeaderTextLength = len(HeaderText)
const LengthTextLength = 5

var connections = map[string]net.Conn{}

func generateData()  {
	//随机生成数据并随机选取Client来发送
	for {
		if len(connections) != 0 {
			text := random.String(8, random.Uppercase)
			var key string
			for key = range connections {
				break
			}
		}
		time.Sleep(time.Duration(5 * int(time.Second)))
	}
}

func main() {
	netListen, err := net.Listen("tcp", "localhost:2048")
	if err != nil {
		log.Println(err)
		return
	}
	defer netListen.Close()

	log.Println("Waiting for clients")
	generateData()
	for {
		conn, err := netListen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println(conn.RemoteAddr().String(), " tcp connect success")
		connections[conn.RemoteAddr().String()] = conn
		go handleConnection(conn)
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

			_, err = conn.Write([]byte("server response"))
			if err != nil {
				log.Println(err)
			}
		}
	}
}
