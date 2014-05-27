package main

import (
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
)

const PORT = ":5000"

type message struct {
	From string `json:"from"`
	Body string `json:"body"`
	DateCreated string `json:"date_created"`
	SID string `json:"sid"`
}

func (m message) MakeUserCommand() UserCommand {
	return UserCommand{ConvertCommand(m.Body), m.From}
}

type TwilioCollector struct {
	listenUrl string
}

type MessageHandler struct {
	messageQueue chan UserCommand
}

func (m MessageHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	msg := message{}
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Print("Error reading in message body")
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(bodyBytes, msg)
	if err != nil {
		log.Print("Error decoding message body")
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	m.messageQueue <- msg.MakeUserCommand()
}

func (t TwilioCollector) Get (queue chan UserCommand) {
	mh := MessageHandler{queue}
	http.Handle(t.listenUrl, mh)

	log.Printf("Listening for Twilio texts on %s\n", PORT)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Print(err.Error())
		log.Fatal("TWILIO http ListenAndServe failed")
	}
}
