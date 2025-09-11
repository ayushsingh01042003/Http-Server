package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Println("Error starting udp server")
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Conn failed", err)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		readData, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Err reading data", err)
		}

		_, err = conn.Write([]byte(readData))
		if err != nil {
			fmt.Println("Err writing to console", err)
		}
	}
}