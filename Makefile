default: build

run:
	GIN_MODE=release go run main.go

build:
	go build -o chat-dingtalk main.go

linux:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o chat-dingtalk main.go

linux-arm:
	CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -o chat-dingtalk main.go

lint:
	env GOGC=25 golangci-lint run --fix -j 8 -v ./...