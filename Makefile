.PHONY: build
build:
	go build -v ./
test:
	go test -v -race ./
bench:
	go test -bench . -benchmem . ./
.DEFAULT_GOAL := build