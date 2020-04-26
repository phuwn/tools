# usage: `make test` or `make test flags="-v"`
test:
	go test $(flags) ./tests/...

example:
	go run example.go