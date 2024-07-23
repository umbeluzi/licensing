.PHONY: build clean run

build:
	go build -o bin/licensing cmd/licensing/main.go

clean:
	rm -rf bin

run:
	go run cmd/licensing/main.go
