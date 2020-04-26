package main

import (
	"github.com/phuwn/tools/errors"
	"github.com/phuwn/tools/log"
)

func errorLog() {
	log.Error(errors.New("error"))
}

func main() {
	log.Info("hey")
	errorLog()
}
