package main

import (
	"log"
	"net/http"
)

const (
	PORT             = ":5000"
	CommandSleepTime = 50
)

func main() {
	commandQueue := make(chan Command, 50)

	// static files
	http.Handle("/", http.FileServer(http.Dir("static/")))
	http.Handle("/index.html", http.FileServer(http.Dir("static/")))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// set up logging
	moveQueue := make(chan Command)
	sh := NewStatHandlerString(moveQueue)
	sh2 := NewStatHandlerJSON(moveQueue)

	// add endpoints
	http.Handle("/stats", sh)
	http.Handle("/stats2", sh2)
	http.Handle("/press", MessageHandler{commandQueue})

	// send off command worker
	go func() {
		for {
			select {
			case cmd := <-commandQueue:
				// send to the old move queue
				moveQueue <- cmd
				log.Printf(cmd.ToString())

				// send off go-routine to execute command
				go EmulatorExecute(ConvertCommand(cmd.Key))
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
