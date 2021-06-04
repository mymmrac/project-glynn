lint-install:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.39.0

lint:
	$(shell go env GOPATH)/bin/golangci-lint run

build-client:
	go build -o bin/glynn cmd/glynn/glynn.go

build-server:
	go build -o bin/glynn-server cmd/glynn-server/glynn-server.go