package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {

	file, err := os.Open("messages.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		con, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		lines := getLinesChannel(con)
		for line := range lines {
			fmt.Printf("%s\n", line)
			// _, err := con.Write([]byte(line))
			// if err != nil {
			// 	fmt.Println("Error writing to connection:", err)
			// }
		}
		con.Close()
		fmt.Println("Connection closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	data := [8]byte{}
	sentence := ""
	line := make(chan string)

	go func() {
		defer close(line)
		defer f.Close()

		for {

			amountRead, err := f.Read(data[:])
			if err != nil {
				if err == io.EOF {
					if sentence != "" {
						//fmt.Printf("read: %s\n", sentence)
						line <- sentence
					} else {
						break
					}
					break
				}
				panic(err)
			}

			str := string(data[:amountRead])
			parts := strings.Split(str, "\n")
			if len(parts) == 1 {
				sentence += parts[0]
				continue
			} else {
				for i := 0; i < len(parts)-1; i++ {
					if parts[i] == "" {
						continue
					}
					if sentence != "" {
						sentence += parts[i] + " "
					}
				}
				if sentence != "" {
					//fmt.Printf("read: %s\n", sentence)
					line <- sentence
					//sentence = ""
				}
			}
			sentence = ""
			sentence += parts[len(parts)-1]
		}
	}()

	return line
}
