package main

import (
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"strings"
)

const PORT = ":5000"

type message struct {
	From string
	Body string
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

func parseTwilio(b string) message {
	kv := make(map[string]string)
	pairs := strings.Split(b, "&")
	for _, pair := range pairs {
		splits := strings.Split(pair, "=")
		if len(splits) == 2 {
			kv[splits[0]] = splits[1]
		}
	}
	if kv["From"] != "" && kv["Body"] != "" {
		return message{kv["From"], kv["Body"]}
	}
	return message{}
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

	msg = parseTwilio(string(bodyBytes))
	if msg.From == "" {
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
