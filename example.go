package main

import (
	"fmt"

	"github.com/phuwn/tools/errors"
	"github.com/phuwn/tools/log"
)

func errorLog() {
	log.Error(errors.New("error"))
}

func errorLogWithDetail() {
	log.Error(errors.Customize(fmt.Errorf("token has expired"), 401, "Unauthorized"))
}

func main() {
	log.Info("hey")
	errorLog()
	errorLogWithDetail()
}
