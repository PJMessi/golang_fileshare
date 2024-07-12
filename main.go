package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	var port string
	flag.StringVar(&port, "port", "", "port number")
	flag.Parse()

	fmt.Println("Press 's' to send files and 'r' to receive files")
	var purpose string
	fmt.Scanln(&purpose)

	if purpose == "s" {
		handleSender(port)
	} else {
		handleReceiver()
	}
}

func handleSender(port string) {
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
		log.Printf("connected to: %s", con.RemoteAddr())

		go connectAndSendFile(con)
	}
}

func connectAndSendFile(con net.Conn) {
	defer con.Close()

	fileName := "testfile.txt"

	// load file
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("err opening file: %s", err)
	}
	defer file.Close()

	// send file
	bytesSent, err := io.Copy(con, file)
	if err != nil {
		log.Fatalf("err sending file: %s", err)
	}

	log.Printf("sent %d bytes", bytesSent)
}

var peers = []string{"localhost:8080"}

func handleReceiver() {
	for _, peer := range peers {
		// connect to sender
		con, err := net.Dial("tcp", peer)
		if err != nil {
			log.Printf("err connecting to peer: %s", err)
			continue
		}

		log.Printf("connected to peer: %s", peer)
		receiveFile(con, "testfile.txt")

		con.Close()
	}
}

func receiveFile(con net.Conn, fileName string) {
	file, err := os.Create("downloaded" + fileName)
	if err != nil {
		log.Fatalf("err creating file: %s", err)
	}
	defer file.Close()

	totalBytes, err := io.Copy(file, con)
	if err != nil {
		log.Fatalf("err receiving file: %s", err)
	}

	log.Printf("downloaded %d bytes", totalBytes)
}
