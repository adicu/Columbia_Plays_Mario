package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)


type Message struct {
	From string `json:"from"`
	Body string `json:"body"`
}

func (m Message) MakeUserCommand() UserCommand {
	return UserCommand{ConvertCommand(m.Body), m.From}
}

type TwilioMessageHandler struct {
	listenUrl string
	messageQueue chan UserCommand
}

func (t TwilioMessageHandler) GetUrl() string {
	return t.listenUrl
}

func parseTwilio(b string) Message {
	kv := make(map[string]string)
	pairs := strings.Split(b, "&")
	for _, pair := range pairs {
		splits := strings.Split(pair, "=")
		if len(splits) == 2 {
			kv[splits[0]] = splits[1]
		}
	}
	if kv["From"] != "" && kv["Body"] != "" {
		return Message{kv["From"], kv["Body"]}
	}
	return Message{}
}

func (m TwilioMessageHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	msg := Message{}
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Print("Error reading in twilio message body")
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	msg = parseTwilio(string(bodyBytes))
	if msg.From == "" {
		log.Print("Error decoding twilio message body")
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	m.messageQueue <- msg.MakeUserCommand()
}
