package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)
	go func() {
		defer f.Close()
		defer close(out)

		buffer := make([]byte, 8)
		line := ""
		for {
			n, err := f.Read(buffer)
			if err == io.EOF {
				break
			}

			i := bytes.IndexByte(buffer, '\n')
			if i == -1 {
				line += string(buffer[:n])
				continue
			}
			line += string(buffer[:i])
			out <- line
			line = string(buffer[i+1:])
		}
	}()

	return out
}

func main() {

	f, err := os.Open("message.txt")
	if err != nil {
		log.Fatal("error: ", err)
	}

	for line := range getLinesChannel(f) {
		fmt.Println("read:", line)
	}
}
