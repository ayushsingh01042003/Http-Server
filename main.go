package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("data.txt")
	if err != nil {
		fmt.Println("Err opening file")
		return
	}

	var output string
	for {
		data := make([]byte, 8)
		_, err := file.Read(data)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Printf("read: %s\n", output)
				break
			}
			fmt.Println("Err reading file")
			return
		}

		if idx := bytes.IndexRune(data, '\n'); idx != -1 {
			output += string(data[:idx])
			data = data[idx + 1:]
			fmt.Printf("read: %s\n", output)
			output = ""
		}
		output += string(data)
	}
}
