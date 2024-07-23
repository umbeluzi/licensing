.PHONY: build clean run

build:
	go build -o bin/licensectl cmd/licensectl/main.go

clean:
	rm -rf bin

run:
	go run cmd/licensectl/main.go
