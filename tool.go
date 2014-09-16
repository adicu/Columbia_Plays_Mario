package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
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

type EmulatorCommand struct {
	Key   string
	Delay time.Duration
}

func ConvertCommand(c string) EmulatorCommand {
	c = strings.ToLower(c)
	switch c {
	// Up
	case "u":
		fallthrough
	case "up":
		return EmulatorCommand{
			Key:   "Up",
			Delay: time.Microsecond * 100,
		}

	// Left
	case "l":
		fallthrough
	case "left":
		return EmulatorCommand{
			Key:   "Left",
			Delay: time.Microsecond * 100,
		}

	// Down
	case "d":
		fallthrough
	case "down":
		return EmulatorCommand{
			Key:   "Down",
			Delay: time.Microsecond * 100,
		}

	// Right
	case "r":
		fallthrough
	case "right":
		return EmulatorCommand{
			Key:   "Right",
			Delay: time.Microsecond * 100,
		}

	// A
	case "a":
		return EmulatorCommand{
			Key:   "z",
			Delay: time.Microsecond * 400,
		}

	// B
	case "b":
		return EmulatorCommand{
			Key:   "x",
			Delay: time.Microsecond * 200,
		}
	}

	return EmulatorCommand{
		Key:   "",
		Delay: 0,
	}
}

func EmulatorExecute(c EmulatorCommand) {
	var keyDown *exec.Cmd
	var keyUp *exec.Cmd
	keyDown = exec.Command(GVBAM_TOOL, "keydown", c.Key)
	keyUp = exec.Command(GVBAM_TOOL, "keyup", c.Key)

	log.Printf("key down %s", c)
	err := keyDown.Run()
	if err != nil {
		log.Println("ERROR: xdotool not functioning properly")
	}

	// sleep the allocated time for this command
	time.Sleep(c.Delay)

	log.Printf("key up %s", c)
	err = keyUp.Run()
	if err != nil {
		log.Println("ERROR: xdotool not functioning properly")
	}
}
