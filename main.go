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

	if purpose == "s" {
		sender.Handle(port)
	} else if purpose == "r" {
		receiver.Handle()
	} else {
		log.Println("invalid input")
	}
}
