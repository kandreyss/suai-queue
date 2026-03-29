BIN_PATH=bin/suai-queue
PATH_TO_MAIN=cmd/bot/main.go

build:
	rm -f $(BIN_PATH)
	go build -o $(BIN_PATH) $(PATH_TO_MAIN)
	strip $(BIN_PATH)

run:
	go run $(PATH_TO_MAIN)

fmt:
	gofmt -w -s -l .
