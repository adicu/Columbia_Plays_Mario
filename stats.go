package main

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	HeldMoves = 40
)

type statHandlerJSON struct {
	newMoves  chan Command
	lastMoves []Command
}

func NewStatHandlerJSON(moves chan Command) *statHandlerJSON {
	sh := statHandlerJSON{moves, make([]Command, HeldMoves)}

	// keeps last 40 commands updated
	go func() {
		for {
			newCommand := <-moves
			copy(sh.lastMoves[1:], sh.lastMoves[:]) // shift by 1
			sh.lastMoves[0] = newCommand
		}
	}()
	return &sh
}

func (s statHandlerJSON) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	respBytes, err := json.Marshal(s.lastMoves)
	if err != nil {
		log.Print("Error encoding stats")
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = resp.Write(respBytes)
	if err != nil {
		log.Print("Error writing stats")
		resp.WriteHeader(http.StatusInternalServerError)
	}
}

type statHandlerString struct {
	newMoves  chan Command
	lastMoves []string
}

func NewStatHandlerString(moves chan Command) *statHandlerString {
	sh := statHandlerString{moves, make([]string, HeldMoves)}

	// keeps last 40 commands updated
	go func() {
		for {
			newCommand := <-moves
			copy(sh.lastMoves[1:], sh.lastMoves[:]) // shift by 1
			sh.lastMoves[0] = newCommand.ToString()
		}
	}()
	return &sh
}

func (s statHandlerString) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	respBytes, err := json.Marshal(s.lastMoves)
	if err != nil {
		log.Print("Error encoding stats")
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = resp.Write(respBytes)
	if err != nil {
		log.Print("Error writing stats")
		resp.WriteHeader(http.StatusInternalServerError)
	}
}
