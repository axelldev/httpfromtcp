package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
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
		if len(line) != 0 {
			out <- line
		}
	}()

	return out
}

func main() {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error: ", err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Connection accepted")

		for line := range getLinesChannel(conn) {
			fmt.Println(line)
		}
		fmt.Println("Connection closed")
	}
}
