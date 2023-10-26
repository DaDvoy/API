BINARY_NAME = api

.PHONY: build all clean

#.DEFAULT_GOAL := build
#SRCS		=  ${shell find ./internal -name "*.go"} main.go

all:	build
		./$(BINARY_NAME)

build:
		go build -o $(BINARY_NAME) ./cmd/main.go

clean:
		rm $(BINARY_NAME)