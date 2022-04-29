.PHONY: build start 

build:
	go build -v ./

start:

	go run ./*.go
	
.DEFAULT_GOAL := start