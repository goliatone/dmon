package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/firstrow/tcp_server"
	"github.com/goliatone/dmon/health"
)

func main() {
	port := "9386"

	if len(os.Args) == 2 {
		port = os.Args[1]
	}

	log.Println("Starting up dmon")

	connectionAddress := fmt.Sprintf("localhost:%s", port)
	server := tcp_server.New(connectionAddress)

	//Handle new clients
	server.OnNewClient(func(c *tcp_server.Client) {
		//New client just connected
		// c.Send("")
		log.Println("New client joined")
	})

	//Handle messages
	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		//new message here!
		message = strings.Trim(message, "\n")
		log.Printf("Command received: \"%s\"", message)

		active, _, message := health.Exec(message)

		log.Println(message)

		if active == false {
			c.Send("KO\n")
			return
		}

		c.Send("OK\n")

	})

	//Handle connection closed
	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		//connection lost
		log.Println("Client left")
	})

	server.Listen()
}
