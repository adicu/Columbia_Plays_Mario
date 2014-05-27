package main

import (
	"bytes"
	"net/http"
	"os/exec"
	"strings"
)

const (
	GVBAM_TOOL = "xdotool"
)

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

	// Start
	case "s":
		fallthrough
	case "start":
		return "Enter"
	}
	return ""
}

// Command to be passed to the emulator
type UserCommand struct {
	key  string
	user string
}

type Message struct {
	From string `json:"from"`
	Body string `json:"body"`
}

func (m Message) MakeUserCommand() UserCommand {
	return UserCommand{ConvertCommand(m.Body), m.From}
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

type GVBAM struct {
	Window string
}

func (g GVBAM) Command(c string) {
	var keyPress *exec.Cmd
	keyPress = exec.Command(GVBAM_TOOL, "key", "--window", g.Window, c)

	var output bytes.Buffer
	keyPress.Stdout = &output
	keyPress.Stderr = &output

	err := keyPress.Run()
	if err != nil {
		panic("xdotool not functioning properly")
	}
}
