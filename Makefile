
EXECUTABLE = pokebot
GO_FILES = main.go \
	tool.go \
	hipchat.go \
	twilio.go \
	stats.go

default: build

build:
	go build -o $(EXECUTABLE) *.go

clean:
	rm $(EXECUTABLE)


