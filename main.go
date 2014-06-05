package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

const (
	PORT            = ":80"
	TwilioEndpoint  = "/twilio"
	HipChatEndpoint = "/hipchat"
	StatsEndpoint   = "/stats"
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
	commandQueue := make(chan UserCommand, 50)

	// set up logging
	moveQueue := make(chan string)
	sh := NewStatHandler(moveQueue)
	http.Handle(StatsEndpoint, sh)

	go getCommands([]CommandCollector{
		TwilioMessageHandler{TwilioEndpoint, commandQueue},
		HipChatCollector{HipChatEndpoint, commandQueue},
	})

	for {
		select {
		case cmd := <-commandQueue:
			moveQueue <- cmd.ToString()
			log.Printf(cmd.ToString())
			emulator.Command(cmd.key)
			time.Sleep(500 * time.Millisecond)
		}
	}
}
