package sender

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type Sender struct {
	chunkSize uint
}

func NewSender(chunkSize uint) *Sender {
	return &Sender{
		chunkSize: chunkSize,
	}
}

func (s *Sender) Handle(port string) error {
	// CREATE A LISTENER
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("err starting listener: %s", err)
	}
	defer listener.Close()
	log.Printf("listening on port: %s", port)

	// LISTEN FOR CLIENTS IN A LOOP
	for {
		con, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("err accepting connection: %s", err)
		}
		log.Printf("connected to receiver: %s", con.RemoteAddr())

		go s.sendFile(con)
	}
}

func (s *Sender) sendFile(con net.Conn) error {
	defer con.Close()

	// REQUEST FILE PATH
	filepath := s.requestFilePath()

	// LOAD THE FILE
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("err opening file: %s", err)
	}
	defer file.Close()

	// SEND FILE NAME SIZE
	if err := s.sendFileNameSize(con, file); err != nil {
		return fmt.Errorf("err sending file name size: %s", err)
	}

	// SEND FILE NAME
	_, err = con.Write([]byte(filepath))
	if err != nil {
		log.Fatalf("err sending filename: %s", err)
	}

	// SEND FILE CONTENT
	if err := s.sendFileContent(con, file); err != nil {
		return fmt.Errorf("err sending file content: %s", err)
	}

	return nil
}

func (s *Sender) sendFileNameSize(con net.Conn, file *os.File) error {
	fileName := file.Name()

	fileNameLen := uint32(len(fileName))

	if err := binary.Write(con, binary.LittleEndian, fileNameLen); err != nil {
		return fmt.Errorf("err sending file name size: %s", err)
	}

	return nil
}

func (s *Sender) sendFileContent(con net.Conn, file *os.File) error {
	chunk := make([]byte, s.chunkSize)

	totalBytesSent := 0
	for {
		// READ A CHUNK
		bytesRead, err := file.Read(chunk)
		if err != nil {
			if err == io.EOF {
				break
			}

			return fmt.Errorf("err reading file chunk: %s", err)
		}

		// SEND THE CHUNK
		// Using con.Write(chunk[:n]) instead of con.Write(chunk) is important
		// because the file.Read(chunk) function doesn’t always fill the buffer
		// completely. It returns the actual number of bytes read, which can be
		// less than the buffer size, especially in the last chunk or if the
		// file is smaller than the buffer size. con.Write(chunk) would send the
		// entire buffer, including any uninitialized or old data, leading to
		// incorrect data transmission.
		_, err = con.Write(chunk[:bytesRead])
		if err != nil {
			return fmt.Errorf("err sending file chunk: %s", err)
		}

		totalBytesSent += bytesRead
	}

	log.Printf("sent %d bytes to receiver", totalBytesSent)

	return nil
}

func (s *Sender) requestFilePath() string {
	fmt.Println("enter the filepath: ")
	var filepath string

	fmt.Scanln(&filepath)

	return filepath
}
