package main

import (
	"os"
	"log"
)

const (
	TwilioEndpoint = "/twilio"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Missing window number")
	}
	emulator := GVBAM{os.Args[1]}
	commandQueue := make(chan UserCommand)

	t := TwilioCollector{TwilioEndpoint}
	go t.Get(commandQueue)

	var cmd UserCommand
	for {
		select {
		case cmd = <-commandQueue:
			log.Printf("Move: %6s By: %s\n", cmd.key, cmd.user)
			emulator.Command(cmd.key)
		}
	}
}
