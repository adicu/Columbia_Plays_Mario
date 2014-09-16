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

	// static files
	http.Handle("/", http.FileServer(http.Dir("static/")))
	http.Handle("/index.html", http.FileServer(http.Dir("static/")))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// set up logging
	moveQueue := make(chan string)
	sh := NewStatHandler(moveQueue)
	http.Handle(StatsEndpoint, sh)

	go getCommands([]CommandCollector{
		TwilioMessageHandler{TwilioEndpoint, commandQueue},
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
