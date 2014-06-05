package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type moveObject struct {
	Move string
}

type statHandler struct {
	newMoves  chan string
	lastMoves []moveObject
}

func NewStatHandler(moves chan string) *statHandler {
	sh := statHandler{moves, make([]moveObject, 20)}
	go func() {
		for {
			newCommand := <-moves
			copy(sh.lastMoves[1:], sh.lastMoves[:])
			sh.lastMoves[0] = moveObject{newCommand}
		}
	}()
	return &sh
}

func (s statHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("Wrong http method, %s, on status endpoint.", req.Method)
		return
	}

	respBytes, err := json.Marshal(s.lastMoves)
	if err != nil {
		log.Print("Error encoding stats")
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(resp, string(respBytes))
}
