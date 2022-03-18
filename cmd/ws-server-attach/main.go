package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Bananenpro/log"
	"github.com/Bananenpro/ws-testserver/attach"
	"github.com/Bananenpro/ws-testserver/cli"
)

func main() {
	log.SetSeverity(log.INFO)

	address := flag.String("address", "127.0.0.1", "The ip address of the server")
	port := flag.Int("port", 80, "The port of the server")
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Printf("Usage: %s <client-id>\n", os.Args[0])
		os.Exit(1)
	}

	client, err := attach.Connect(fmt.Sprintf("ws://%s:%d/attach/%s", *address, *port, flag.Arg(0)))
	if err != nil {
		log.Fatal("Failed to connect to server:", err)
	}

	log.Info("Connection established.")

	go client.Listen()

	for {
		msg := cli.AskForMessage()
		if msg != "" {
			err := client.Send(msg)
			if err != nil {
				log.Error("Failed to send message to server:", err)
			}
		}
	}
}
