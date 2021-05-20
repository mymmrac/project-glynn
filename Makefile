lint-install:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.39.0

lint:
	$(shell go env GOPATH)/bin/golangci-lint run

test:
	go test -coverprofile=coverprofile.out ./...

validate: lint test

mock-install:
	go install github.com/golang/mock/mockgen@v1.5.0

mockgen:
	mockgen -destination internal/mocks/repository.go -package mocks github.com/mymmrac/project-glynn/pkg/repository Repository
