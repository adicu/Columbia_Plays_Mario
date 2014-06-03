package main

import (
	"log"
	"net/http"
	"os"
)

const (
	PORT            = ":5000"
	TwilioEndpoint  = "/twilio"
	HipChatEndpoint = "/hipchat"
)

func checkEnvVariables() {
	// check CLI args
	if len(os.Args) != 2 {
		log.Fatal("Missing window number")
	}
}

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
	checkEnvVariables()
	emulator := GVBAM{os.Args[1]}
	commandQueue := make(chan UserCommand)

	go getCommands([]CommandCollector{
		TwilioMessageHandler{TwilioEndpoint, commandQueue},
		HipChatCollector{HipChatEndpoint, commandQueue},
	})

	for {
		select {
		case cmd := <-commandQueue:
			log.Printf("Move: %6s Via %10s By: %s\n", cmd.key, cmd.via, cmd.user)
			emulator.Command(cmd.key)
		}
	}
}
