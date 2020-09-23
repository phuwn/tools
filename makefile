# usage: `make test` or `make test flags="-v"`
test:
	go test $(flags) ./...

example:
	go run example.go