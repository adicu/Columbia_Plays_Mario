package main

import (
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
		log.Printf("Wrong http method, %s, on Twilio endpoint.", req.Method)
		return
	}

	msg := Message{req.FormValue("From"), req.FormValue("Body")}
	if msg.From == "" || msg.Body == "" {
		log.Print("Error decoding twilio message body")
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	m.messageQueue <- msg.MakeUserCommand()
}
