package errors

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/phuwn/tools/util"
)

var (
	// RecordNotFound error message
	RecordNotFound = "record not found"
)

// CError - Custom error for debugging and json responsing
type CError struct {
	Code    int      `json:"-"`
	Message string   `json:"error"`
	Details []string `json:"details,omitempty"`
	stack   *util.Stack
}

func (e CError) Error() string {
	return e.Message
}

func (e CError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, e.Message)
			if len(e.Details) > 0 {
				io.WriteString(s, "\nDetails:\n- "+strings.Join(e.Details, "\n- "))
			}
			e.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.Message)
	case 'q':
		fmt.Fprintf(s, "%q", e.Message)
	}
}

// New - generate a new custom error with input info and internal server status as default status
func New(msg string, args ...interface{}) error {
	return &CError{
		Code:    http.StatusInternalServerError,
		Message: fmt.Sprintf(msg, args...),
		stack:   util.Callers(),
	}
}

// Customize - customize normal error to CError, push new message ahead of the details stack
func Customize(code int, msg string, err error) error {
	if err == nil {
		return nil
	}

	ce, ok := err.(*CError)
	if ok {
		if code != 0 {
			ce.Code = code
		}
		if msg != "" {
			ce.Details = append([]string{ce.Message}, ce.Details...)
			ce.Message = msg
		}
		return ce
	}

	ce = &CError{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
		stack:   util.Callers(),
	}

	if code != 0 {
		ce.Code = code
	}
	if msg != "" {
		ce.Details = []string{err.Error()}
		ce.Message = msg
	}
	return ce
}

// IsRecordNotFound - check if an error is a RecordNotFound error
func IsRecordNotFound(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == RecordNotFound
}
