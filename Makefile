run:
	go run ./main.go

build:
	go build -o bin/app ./main.go

test:
	go test ./...
