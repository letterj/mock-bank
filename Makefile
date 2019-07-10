# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=fc2-mock-bank
BINARY_MACOS=$(BINARY_NAME).macos
BINARY_LINUX=$(BINARY_NAME).linux
FC2BANK=bank.db


all: clean build

build: 
	$(GOBUILD) -o $(BINARY_MACOS) -v

test: 
	$(GOTEST) -v ./...
		
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_MACOS)
	rm -f $(BINARY_LINUX)
	rm -f $(FC2BANK)

run:
	$(GOBUILD) -o $(BINARY_MACOS) -v ./...
	./$(BINARY_MACOS) 

deps:
	$(GOGET) github.com/mattn/go-sqlite3
	$(GOGET) github.com/gorilla/mux
	$(GOGET) github.com/gorilla/handlers

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_LINUX) -v
    
# docker-build:
	# docker run --rm -it -v "$(GOPATH)":/go -w $HOME/go/src/github.com/9thGear/fc2-mock-bank golang:latest go build -o "$(BINARY_LINUX" -v
