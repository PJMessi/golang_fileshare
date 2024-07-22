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

type Receiver struct {
	chunkSize uint
}

func NewReceiver(chunkSize uint) *Receiver {
	return &Receiver{
		chunkSize: chunkSize,
	}
}

func (r *Receiver) Handle() error {
	var peers = []string{"localhost:8080"}

	for _, peer := range peers {
		// CONNECT TO SENDER
		con, err := net.Dial("tcp", peer)
		if err != nil {
			log.Printf("err connecting to peer: %s", err)
			continue
		}

		log.Printf("connected to peer: %s", peer)

		// RECEIVE FILE FROM SENDER
		if err = r.receiveFile(con); err != nil {
			return fmt.Errorf("err receiving file: %s", err)
		}

		if err = con.Close(); err != nil {
			return fmt.Errorf("err closing connection: %s", err)
		}
	}

	return nil
}

func (r *Receiver) receiveFile(con net.Conn) error {
	// RECEIVE FILE NAME
	filePath, err := r.receiveFileName(con)
	if err != nil {
		return fmt.Errorf("err receiving file name: %s", err)
	}

	// PREPARE PATH TO SAVE THE FILE
	destFilePath := r.prepareDestFilePath(filePath)

	// CREATE FILE
	file, err := os.Create(destFilePath)
	if err != nil {
		return fmt.Errorf("err creating dest file: %s", err)
	}
	defer file.Close()

	// SAVE CONTENT TO THE FILE
	if err = r.receiveAndSaveFileContent(con, file); err != nil {
		return fmt.Errorf("err receiving and saving file content: %s", err)
	}

	return nil
}

func (r *Receiver) receiveAndSaveFileContent(con net.Conn, file *os.File) error {
	chunk := make([]byte, r.chunkSize)

	totalBytesReceived := 0

	for {
		bytesRead, err := con.Read(chunk)
		if err != nil {
			if err == io.EOF {
				break
			}

			return fmt.Errorf("err receiving file chunk: %s", err)
		}

		totalBytesReceived += bytesRead

		_, err = file.Write(chunk[:bytesRead])
		if err != nil {
			return fmt.Errorf("err writing chunk to the file: %s", err)
		}
	}

	return nil
}

func (r *Receiver) receiveFileName(con net.Conn) (string, error) {
	fileNameLen, err := r.receiveFileNameLen(con)
	if err != nil {
		return "", fmt.Errorf("err receiving file name len: %s", err)
	}

	nameBuf := make([]byte, fileNameLen)
	_, err = io.ReadFull(con, nameBuf)
	if err != nil {
		return "", fmt.Errorf("err receiving file name: %s", err)
	}

	return string(nameBuf), nil
}

func (r *Receiver) receiveFileNameLen(con net.Conn) (uint32, error) {
	// INFO: match tye type with the sender
	var uintType uint32

	lenBuf := make([]byte, unsafe.Sizeof(uintType))

	_, err := io.ReadFull(con, lenBuf)
	if err != nil {
		return 0, fmt.Errorf("err receiving file name length: %s", err)
	}

	fileNameLen := binary.LittleEndian.Uint32(lenBuf)
	return fileNameLen, nil
}

func (r *Receiver) prepareDestFilePath(filePath string) string {
	fileExt := path.Ext(filePath)

	destFilePath := fmt.Sprintf("%d%s", time.Now().Unix(), fileExt)

	return destFilePath
}
