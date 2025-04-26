package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/l2thet/httpfromtcp/internal/request"
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
		httpData, err := request.RequestFromReader(con)
		if err != nil {
			fmt.Println("Error reading request:", err)
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", httpData.RequestLine.Method)
		fmt.Printf("- Target: %s\n", httpData.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", httpData.RequestLine.HttpVersion)

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
