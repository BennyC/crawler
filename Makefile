test:
	go test ./...
coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out
build:
	go build cmd/crawler/main.go
