package main

import (
	"fmt"
	"github.com/fragmenta/server"
	"github.com/techfront/core/src"
)

func main() {

	// If we have no config, bootstrap first by generating config/migrations
	if src.RequiresBootStrap() {
		err := src.Bootstrap()
		if err != nil {
			fmt.Printf("Error bootstrapping server %s\n", err)
			return
		}
	}

	// Setup server
	server, err := server.New()
	if err != nil {
		fmt.Printf("Error creating server %s", err)
		return
	}

	src.Setup(server)

	// Inform user of server setup
	server.Logf("#info Starting server in %s mode on port %d", server.Mode(), server.Port())

	// Start the server
	err = server.Start()
	if err != nil {
		server.Fatalf("Error starting server %s", err)
	}

}
