package main

import (
	"flag"

	"github.com/Bananenpro/log"
	"github.com/Bananenpro/ws-testserver/server"
)

func main() {
	port := flag.Int("port", 80, "The port to listen on")
	flag.Parse()

	log.SetSeverity(log.TRACE)

	server := server.New()
	server.Listen(*port)
}
