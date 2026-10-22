package main

import (
	"crache/cache"

	"net"
	"log"
	"bufio"
	"strconv"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatal("Error listening: ", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	log.Println("connection created")

	reader := bufio.NewReader(conn)

	msgLength, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Message Length Read error: %v", err)
		return
	}
	len, err := strconv.Atoi(string(msgLength)[1:2])
	if err != nil {
		log.Printf("Error converting string to int: %v", err)
		return
	}

	log.Printf("length of message is %d ints", len)

	buf := make([]string, 0, 1024)

	for i := 0; i < len; i++ {
		_, err := reader.ReadString('\n') // length of recvStr
		if err != nil {
			log.Printf("Message Reading Error: %v", err)
			return
		}

		recvStr, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Message Reading Error: %v", err)
			return
		}
		buf = append(buf, strings.TrimSpace(recvStr))
	}

	// switch case for resp command
	switch buf[0] {
	case "COMMAND":
		if buf[1] == "DOCS" {
			log.Printf("return summary")
		} else {
			log.Printf("unknown command")
		}
	case "GET":
		cache.HandleGet()
		log.Printf("get statement")
	case "SET": 
		cache.HandleSet()
		log.Printf("set statement")
	}

	log.Printf("Hi!")
	log.Printf("%v", buf)

}