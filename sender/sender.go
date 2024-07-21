package sender

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func Handle(port string) {
	// listen for clients
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("err starting listener: %s", err)
	}
	defer listener.Close()
	log.Printf("listening on port: %s", port)

	// keep listening for clients
	for {
		con, err := listener.Accept()
		if err != nil {
			log.Fatalf("err accepting connection: %s", err)
		}
		log.Printf("connected to receiver: %s", con.RemoteAddr())

		go sendFile(con)
	}
}

func sendFile(con net.Conn) {
	defer con.Close()

	filepath := requestFilePath()
	log.Println("filepath: ", filepath)

	// load file
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("err opening file: %s", err)
	}
	defer file.Close()

	// sending file name first
	log.Println("the filename to be sent")
	_, err = con.Write([]byte(filepath))
	if err != nil {
		log.Fatalf("err sending filename: %s", err)
	}

	// creating buffer to hold 1024 bytes (1 kb)
	chunk := make([]byte, 1024)

	totalBytesSent := 0
	for {
		// read a chunk
		n, err := file.Read(chunk)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("err reading a chunk: %s", err)
		}

		// send the chunk
		// Using con.Write(chunk[:n]) instead of con.Write(chunk) is important
		// because the file.Read(chunk) function doesn’t always fill the buffer
		// completely. It returns the actual number of bytes read, which can be
		// less than the buffer size, especially in the last chunk or if the
		// file is smaller than the buffer size. con.Write(chunk) would send the
		// entire buffer, including any uninitialized or old data, leading to
		// incorrect data transmission.
		_, err = con.Write(chunk[:n])
		if err != nil {
			log.Printf("err sending chunk: %s", err)
		}

		totalBytesSent += n
	}

	log.Printf("sent %d bytes", totalBytesSent)
}

func requestFilePath() string {
	fmt.Println("enter the filepath: ")
	var filepath string
	fmt.Scanln(&filepath)
	return filepath
}
