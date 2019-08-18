VERSION := $(shell git tag | tail -1)
BUILD := $(shell git log --pretty=%h | head -1)
APP_NAME := emmstats

MAKEFLAGS += --silent

build:
	# MacOS binary
	go build -o bin/$(APP_NAME)-osx -ldflags "-X main.version=$(VERSION) -X main.build=$(BUILD) -X main.toolName=$(APP_NAME)" 
	gzip -f bin/$(APP_NAME)-osx 
	
	# Linux Binary
	GOOS=linux GOARCH=386 go build -o bin/$(APP_NAME)-linux -ldflags "-X main.version=$(VERSION) -X main.build=$(BUILD) -X main.toolName=$(APP_NAME)" 
	gzip -f bin/$(APP_NAME)-linux 
	
	# Windows Binary
	# Dependency required for Windows build only	
	go get github.com/konsorten/go-windows-terminal-sequences
	GOOS=windows GOARCH=386 go build -o bin/$(APP_NAME)-win -ldflags "-X main.version=$(VERSION) -X main.build=$(BUILD) -X main.toolName=$(APP_NAME)" 
	gzip -f bin/$(APP_NAME)-win 

clean:
	rm -rvf bin

test:
	echo $(VERSION)
	echo $(BUILD)
