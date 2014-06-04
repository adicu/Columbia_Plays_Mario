package main

import (
	"fmt"
	"log"
	"net/http"
)

type Message struct {
	From string `json:"from"`
	Body string `json:"body"`
}

type TwilioMessageHandler struct {
	listenUrl    string
	messageQueue chan UserCommand
}

func (t TwilioMessageHandler) GetUrl() string {
	return t.listenUrl
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

	cmd := UserCommand{ConvertCommand(msg.Body), "Twilio", msg.From}
	if cmd.key == "" {
		log.Printf("Invalid gameboy move, \"%s\"\n", msg.Body)
		resp.WriteHeader(http.StatusBadRequest)
		resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := fmt.Fprint(resp, "Invalid move, please use:\na / b / l(eft) / u(p) / r(ight) / d(own) / start / select")
		if err != nil {
			log.Printf("Error while writing http response, %s\n", err.Error())
		}
		return
	}

	m.messageQueue <- cmd
}
