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
	server.OnNewMessage(func(c *tcp_server.Client, payload string) {
		//new message here!
		payload = strings.Trim(payload, "\n")
		log.Printf("Command received: \"%s\"", payload)

		response := health.Exec(payload)

		log.Println(response.Message)

		if response.Success == false {
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
