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

	receiver := receiver.NewReceiver(1024)

	if purpose == "s" {
		sender.Handle(port)

	} else if purpose == "r" {
		if err := receiver.Handle(); err != nil {
			log.Fatalf("err receiving file from the sender: %s", err)
		}

	} else {
		log.Println("invalid input")
	}
}
