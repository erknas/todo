
run: build
	@./bin/app
build:
	@go build -o bin/app cmd/main.go
test:
	@go test ./tests
