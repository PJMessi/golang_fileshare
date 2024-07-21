package receiver

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"time"
	"unsafe"
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
	// READ FILE NAME LENGTH
	var uintType uint32 // match tye type with the sender
	lenBuf := make([]byte, unsafe.Sizeof(uintType))
	_, err := io.ReadFull(con, lenBuf)
	if err != nil {
		log.Fatalf("err receiving file name length: %s", err)
	}
	fileNameLen := binary.LittleEndian.Uint16(lenBuf)
	log.Println("file name length: ", fileNameLen)

	// READ FILE NAME
	nameBuf := make([]byte, fileNameLen)
	_, err = io.ReadFull(con, nameBuf)
	if err != nil {
		log.Fatalf("err receiving file name: %s", err)
	}
	sourceFilePath := string(nameBuf)
	log.Println("file name: ", sourceFilePath)

	// PREPARE FILE NAME FOR SAVING
	fileName := prepareFileName(sourceFilePath)

	// CREATING A FILE TO PUT FILE CONTENT
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("err creating file: %s", err)
	}
	defer file.Close()

	// creating buffer to hold 1024 bytes (1 kb)
	chunk := make([]byte, 1024)
	totalBytesReceived := 0

	for {
		n, err := con.Read(chunk)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("err receiving chunk: %s", err)
			if err := os.Remove(file.Name()); err != nil {
				log.Fatalf("err removing file: %s", err)
			}
		}

		totalBytesReceived += n

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

func prepareFileName(filePath string) string {
	fileExt := path.Ext(filePath)

	destFileName := fmt.Sprintf("%d%s", time.Now().Unix(), fileExt)

	log.Println("filepath: ", filePath)
	log.Println("filename: ", fileExt)
	log.Println("destFileName: ", destFileName)

	return destFileName
}
