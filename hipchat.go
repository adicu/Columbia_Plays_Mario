package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
)

type HipChatCollector struct {
	listenUrl string
	messageQueue chan UserCommand
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

	msg := Message{}
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Print("Error reading in hipchat message body")
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(bodyBytes, &msg)
	if err != nil {
		log.Print("Error parsing hipchat json body")
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.messageQueue <- msg.MakeUserCommand()
}
