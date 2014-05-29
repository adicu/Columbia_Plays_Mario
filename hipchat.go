package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type HipChatCollector struct {
	listenUrl    string
	messageQueue chan UserCommand
}

type HipChatHook struct {
	Event string `json:"event"`
	Item  struct {
		Message struct {
			Message string `json:"message"`
			Date    string `json:"date"`
			From    struct {
				Name string `json:"name"`
			} `json:"from"`
		} `json:"message"`
	} `json:"item"`
	WebhookId string `json:"webhook_id"`
}

func (m HipChatHook) MakeUserCommand() UserCommand {
	return UserCommand{ConvertCommand(m.Item.Message.Message), m.Item.Message.From.Name}
}

func (h HipChatCollector) GetUrl() string {
	return h.listenUrl
}

func (h HipChatCollector) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("Wrong http method, %s, on HipChat endpoint.", req.Method)
		return
	}

	msg := HipChatHook{}
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Print("Error reading in hipchat message body")
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(bodyBytes, &msg)
	if err != nil {
		log.Print("Error parsing hipchat webhook json body")
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.messageQueue <- msg.MakeUserCommand()
}
