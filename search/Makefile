
pre:
	go mod tidy

unit:
	go test ./... -cover

build: pre
	mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -gcflags "-mmm" . -o bin/server
