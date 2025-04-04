package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {

	file, err := os.Open("messages.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	lines := getLinesChannel(file)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
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
