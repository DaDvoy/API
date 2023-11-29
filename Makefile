BINARY_NAME = api

.PHONY: build all clean

export CONFIG_PATH = ./config/local.yaml

all:	build
		./$(BINARY_NAME)

build:
		go build -o $(BINARY_NAME) ./cmd/main.go

clean:
		rm $(BINARY_NAME)

docker:
		docker run --name some-postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d postgres:latest