package main

import (
	"log"
	"net/http"
	"time"
)

const (
	PORT             = ":5000"
	Endpoint         = "/press"
	StatsEndpoint    = "/stats"
	CommandSleepTime = 50
)

func main() {
	commandQueue := make(chan Command, 50)

	// static files
	http.Handle("/", http.FileServer(http.Dir("static/")))
	http.Handle("/index.html", http.FileServer(http.Dir("static/")))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// set up logging
	moveQueue := make(chan string)
	sh := NewStatHandler(moveQueue)

	// add endpoints
	http.Handle(StatsEndpoint, sh)
	http.Handle("/press", MessageHandler{commandQueue})

	// send off command worker
	go func() {
		for {
			select {
			case cmd := <-commandQueue:
				moveQueue <- cmd.ToString()
				log.Printf(cmd.ToString())
				EmulatorCommand(cmd.Key)
				time.Sleep(CommandSleepTime * time.Millisecond)
			}
		}
	}()

	// start webserver
	log.Printf("Listening for commands on %s\n", PORT)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Print(err.Error())
		log.Fatal("HTTP ListenAndServe failed")
	}
}
