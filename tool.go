package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

const (
	GVBAM_TOOL = "xdotool"
)

type Command struct {
	Username string `json:"username,omitempty"`
	Key      string `json:"key"`
}

func (cmd Command) ToString() string {
	return fmt.Sprintf("Move: %6s By: %s\n", cmd.Key, cmd.Username)
}

// Objects that collect commands to be passed to the emulator
type CommandCollector interface {
	GetUrl() string
	http.Handler
}

// Represents an emulator
type Emulator interface {
	Command(string)
}

// MessageHandler wraps the command pipeline
type MessageHandler struct {
	messageQueue chan Command
}

func (m MessageHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// valid http method
	if req.Method != "POST" {
		log.Println("Wrong HTTP header")
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// read in json to buffer
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("Failed to read in request body")
		resp.WriteHeader(400)
		return
	}

	// parse buffer into a command
	var msg Command
	err = json.Unmarshal(bodyBytes, &msg)
	if err != nil {
		log.Println(string(bodyBytes))
		log.Printf("Failed to unmarshal JSON => %s", err.Error())
		resp.WriteHeader(400)
		return
	}

	// validate the command
	if msg.Key == "" {
		log.Printf("Invalid gameboy move, \"%s\"\n", msg.Key)
		resp.Header().Set("Content-Type", "text/plain")
		_, err := fmt.Fprint(resp, "Invalid move, please use:\na / b / l(eft) / u(p) / r(ight) / d(own) / start / select")
		if err != nil {
			log.Printf("Error while writing http response, %s\n", err.Error())
		}
		return
	}

	// send the command into the pipeline
	go func() {
		m.messageQueue <- msg
	}()
}

func ConvertCommand(c string) string {
	c = strings.ToLower(c)
	switch c {
	// Up
	case "u":
		fallthrough
	case "up":
		return "Up"

	// Left
	case "l":
		fallthrough
	case "left":
		return "Left"

	// Down
	case "d":
		fallthrough
	case "down":
		return "Down"

	// Right
	case "r":
		fallthrough
	case "right":
		return "Right"

	// A
	case "a":
		return "z"

	// B
	case "b":
		return "x"
	}

	return ""
}

func EmulatorCommand(c string) {
	var keyPress *exec.Cmd
	keyPress = exec.Command(GVBAM_TOOL, "key", c);

	var output bytes.Buffer
	keyPress.Stdout = &output
	keyPress.Stderr = &output

	err := keyPress.Start()
	if err != nil {
		log.Println("ERROR: xdotool not functioning properly")
	}
}
