package main

import (
	"bytes"
	"fmt"
	"io"
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

		for line := range getLinesChannel(conn) {
			fmt.Printf("read: %s\n", line)
		}
		conn.Close()
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	strChan := make(chan string)

	var output string
	go func() {
		defer f.Close()
		defer close(strChan)
		
		for {
			data := make([]byte, 8)
			_, err := f.Read(data) // blocking untill data is send

			if err != nil {
				if err.Error() == "EOF" {
					strChan <- output
					break
				}
				fmt.Printf("Error reading the file: %s", err.Error())
				return
			}

			if idx := bytes.IndexRune(data, '\n'); idx != -1 {
				output += string(data[:idx])
				data = data[idx + 1:]
				strChan <- output
				output = ""
			}
			output += string(data)
		}
	}()

	return strChan
}