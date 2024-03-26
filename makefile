rs:
	go run ./cmd/server/main.go

build:
	go build -o ./bin/redis/server ./cmd/server/main.go

PHONY: rs build



