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

	connectionAddress := fmt.Sprintf("0.0.0.0:%s", port)
	server := tcp_server.New(connectionAddress)

	server.OnNewClient(func(c *tcp_server.Client) {
		log.Println("New client joined")
	})

	server.OnNewMessage(func(c *tcp_server.Client, payload string) {
		payload = strings.Trim(payload, "\n")

		log.Printf("Command received: \"%s\"", payload)

		//A) Parse payload.
		//B) Get handler from payload

		//C) Execute handler
		response, err := health.Exec(payload)

		if err != nil {
			log.Printf("Error executing health check: %s", err)
			c.Send("KO\n")
			return
		}

		log.Println(response.Message)

		if response.Success {
			c.Send("OK\n")
			return
		}

		c.Send("KO\n")
	})

	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		log.Println("Client left")
	})

	server.Listen()
}
