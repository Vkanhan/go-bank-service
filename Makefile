build: 
	@go build -o bin/go-bank-service

run: build
	@./bin/go-bank-service

test:
	@go test -v ./...

