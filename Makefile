BIN=bin/suai-queue
PATH_TO_MAIN=cmd/bot/main.go

build:
	rm -f bin/suai-queue
	go build -o $(BIN) $(PATH_TO_MAIN)

run:
	go run $(PATH_TO_MAIN)

fmt:
	gofmt -w -s -l .