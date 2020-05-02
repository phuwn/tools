# usage: `make test` or `make test flags="-v"`
test:
	go test $(flags) ./test/...

example:
	go run example.go