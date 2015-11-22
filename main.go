package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/firstrow/tcp_server"
	"github.com/goliatone/dmon/data"
	"github.com/goliatone/dmon/health"
)

func init() {
	//Retrieve the check bash script and restore it
	err := data.RestoreAsset(".", "bin/check")

	if err != nil {
		log.Println("Error creating asset", err)
		os.Exit(1)
	}
}

func main() {
	server := tcp_server.New("localhost:9386")
	log.Println("Starting up dmon on port 9386")

	server.OnNewClient(func(c *tcp_server.Client) {
		//New client just connected
		// c.Send("")
	})

	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		//new message here!
		message = strings.Trim(message, "\n")
		log.Println(message)

		active, _, message := health.Exec(message)
		log.Println(message)

		if active == false {
			c.Send("KO\n")
			return
		}

		c.Send("OK\n")

	})

	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		//connection lost
		log.Println("Client left")
	})

	server.Listen()
}

func getPwd() string {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	return pwd
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
