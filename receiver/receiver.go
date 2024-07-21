package receiver

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"time"
)

func Handle() {
	var peers = []string{"localhost:8080"}

	for _, peer := range peers {
		// connect to sender
		con, err := net.Dial("tcp", peer)
		if err != nil {
			log.Printf("err connecting to peer: %s", err)
			continue
		}

		log.Printf("connected to peer: %s", peer)
		receiveFile(con)

		con.Close()
	}
}

func receiveFile(con net.Conn) {
	// creating buffer to hold 1024 bytes (1 kb)
	chunk := make([]byte, 1024)
	totalBytesReceived := 0

	fileNameProcessed := false
	var filename string
	var file *os.File
	for {
		n, err := con.Read(chunk)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("err receiving chunk: %s", err)
			if filename != "" {
				if err := os.Remove(filename); err != nil {
					log.Fatalf("err removing file: %s", err)
				}
			}
		}

		totalBytesReceived += n

		// first chunk is the file name, so extract fileName from the first chunk
		if !fileNameProcessed {
			filename := prepareFileName(chunk[:n])
			file, err = os.Create(filename)
			if err != nil {
				log.Fatalf("err creating file: %s", err)
			}
			defer file.Close()
			fileNameProcessed = true
			continue
		}

		_, err = file.Write(chunk[:n])
		if err != nil {
			fmt.Printf("err writing chunk to a file: %s", err)
			if err := os.Remove(file.Name()); err != nil {
				log.Fatalf("err removing file: %s", err)
			}
		}
	}

	log.Printf("received %d bytes", totalBytesReceived)
}

func prepareFileName(fileNameChunk []byte) string {
	filePath := string(fileNameChunk)
	fileExt := path.Ext(filePath)

	destFileName := fmt.Sprintf("%d%s", time.Now().Unix(), fileExt)

	log.Println("filepath: ", filePath)
	log.Println("filename: ", fileExt)
	log.Println("destFileName: ", destFileName)

	return destFileName
}
