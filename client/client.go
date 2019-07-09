package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
)

var waitGroup sync.WaitGroup

func disconnect(conn net.Conn)  {
	conn.Close()
	waitGroup.Done()
}

func generate(infoData chan string)  {
	index := 0
	for {
		index++
		if rand.Intn(200) == 199 {
			//退出程序
			infoData <- "quit"
			return
		} else {
			dic := make(map[string]interface{})
			dic["index"] = index
			dic["timestamp"] = time.Now().Format(time.RFC3339)
			jsonString, err := json.Marshal(dic)
			if err != nil {
				log.Println(err)
			}
			infoData <- string(jsonString)

			time.Sleep(time.Duration(rand.Intn(5) * int(time.Second)))
		}
	}
}

func send(conn net.Conn, infoData chan string) {
	for {
		jsonString := <- infoData
		if jsonString == "quit" {
			disconnect(conn)
			return
		}

		//按协议发送数据
		length := len(jsonString)
		if length > 99999 {
			//Header中标识字符串长度的最大为99999
			panic("data is too long to send")
		}
		lengthText := strconv.Itoa(length)
		textLength := fmt.Sprintf("%05s", lengthText)[:5]
		headerText := append([]byte("Header"), textLength...)
		jsonString = string(append(headerText, jsonString...))
		_, err := conn.Write([]byte(jsonString))
		if err != nil {
			log.Println(err)
		}
		log.Println("send : ", string(jsonString))
	}
}

func read(conn net.Conn) {
	buffer := make([]byte, 2048)
	for {
		log.Println("waiting for read...")
		reqLen, err := conn.Read(buffer)
		if err != nil {
			log.Println(err)
			disconnect(conn)
			return
		}
		log.Println("read length : ", reqLen)
		received := string(buffer[:reqLen-1])
		log.Println(received)

		//服务端请求关闭
		if received == "quit" {
			disconnect(conn)
			return
		}
	}
}

func main()  {
	server := "127.0.0.1:2048"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		log.Println(err)
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Println(err)
		return
	}

	waitGroup.Add(1)

	log.Println("connect success")

	infoData := make(chan string)
	go send(conn, infoData)
	go read(conn)
	go generate(infoData)  //模拟业务逻辑

	waitGroup.Wait()
}