package main

import (
	"log"
	"net/http"
	"os"
)

const (
	PORT            = ":80"
	TwilioEndpoint  = "/twilio"
	HipChatEndpoint = "/hipchat"
)

func getCommands(handlers []CommandCollector) {
	for _, h := range handlers {
		http.Handle(h.GetUrl(), h)
	}

	log.Printf("Listening for commands on %s\n", PORT)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Print(err.Error())
		log.Fatal("HTTP ListenAndServe failed")
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Missing window number")
	}
	emulator := GVBAM{os.Args[1]}
	commandQueue := make(chan UserCommand)

	go getCommands([]CommandCollector{
		TwilioMessageHandler{TwilioEndpoint, commandQueue},
		HipChatCollector{HipChatEndpoint, commandQueue},
	})

	for {
		select {
		case cmd := <-commandQueue:
			log.Printf("Move: %6s By: %s\n", cmd.key, cmd.user)
			emulator.Command(cmd.key)
		}
	}
}
