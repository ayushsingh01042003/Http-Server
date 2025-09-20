package main

import (
	"fmt"
	"httpserver/internal/request"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println("Cannot start the server, ERR-", err.Error())
		return
	}

	for {
		conn, err := listener.Accept() // blocks till connections request send
		if err != nil {
			fmt.Println("Cannot accept connection, ERR-", err)
		}
		fmt.Println("Client connected")

		go func(conn net.Conn) {
			reqData, err := request.RequestFromReader(conn)
			if err != nil {
				fmt.Println("Error reading request", err)
				return 
			}
			reqData.RequestLine.PrintData()
		} (conn)
	}
}
