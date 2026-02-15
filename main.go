package main

import (
	"log"
	"polyserver/signaling"
)

func main() {
	log.Println("Starting...")

	server := signaling.NewServer()

	if err := server.Connect(); err != nil {
		log.Fatal(err)
	}

	go server.Start()

	if err := server.CreateInvite(); err != nil {
		log.Fatal(err)
	}

	select {} // keep program alive
}
