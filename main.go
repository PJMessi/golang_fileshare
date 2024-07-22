package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/pjmessi/go_file_share/receiver"
	"github.com/pjmessi/go_file_share/sender"
)

func main() {
	var port string
	flag.StringVar(&port, "port", "", "port number")
	flag.Parse()

	fmt.Println("Press 's' to send files and 'r' to receive files")
	var purpose string
	fmt.Scanln(&purpose)

	udpDiscoveryPort := uint(9999)
	chunkSize := uint(1024)
	receiver := receiver.NewReceiver(chunkSize, udpDiscoveryPort)
	sender := sender.NewSender(chunkSize, udpDiscoveryPort)

	if purpose == "s" {
		if err := sender.Handle(port); err != nil {
			log.Fatalf("err starting sender: %s", err)
		}

	} else if purpose == "r" {
		if err := receiver.Handle(); err != nil {
			log.Fatalf("err receiving file from the sender: %s", err)
		}

	} else {
		log.Println("invalid input")
	}
}
