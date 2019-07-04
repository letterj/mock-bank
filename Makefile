# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=fcfcmockbank
BINARY_UNIX=$(BINARY_NAME)_unix
FC2BANK=bank.db


all: clean build

build: 
	$(GOBUILD) -o $(BINARY_NAME) -v

test: 
	mkdir -p $(DB_LOCATION)
	$(GOTEST) -v ./...
		
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(FC2BAND)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME) default

deps:
	$(GOGET) github.com/mattn/go-sqlite3
	$(GOGET) github.com/gorilla/mux
	$(GOGET) github.com/gorilla/handlers

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
    
docker-build:
	docker run --rm -it -v "$(GOPATH)":/go -w $HOME/go/src/github.com/letterj/fcfcbank golang:latest go build -o "$(BINARY_UNIX)" -v