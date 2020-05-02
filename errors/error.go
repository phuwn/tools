package errors

import (
	"errors"
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

// NewResp - generate a new custom error with provided info
func NewResp(code int, msg string, details ...string) error {
	return &CError{
		Code:    code,
		Message: msg,
		Details: details,
		stack:   util.Callers(),
	}
}

// UpdateCode - update status of an error
func UpdateCode(err error, code int) error {
	if err == nil {
		return nil
	}
	ce, ok := err.(*CError)
	if ok {
		ce.Code = code
		return ce
	}

	return &CError{
		Code:    code,
		Message: err.Error(),
		stack:   util.Callers(),
	}
}

// AddDetails - push new details to error's details stack
func AddDetails(err error, details ...string) error {
	if err == nil {
		return nil
	}
	ce, ok := err.(*CError)
	if ok {
		ce.Details = append(details, ce.Details...)
		return ce
	}

	return &CError{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
		Details: details,
		stack:   util.Callers(),
	}
}

// Overload - push new message into current error, the old message will be ahead of the details stack
func Overload(msg string, err error) error {
	if err == nil {
		return errors.New(msg)
	}
	ce, ok := err.(*CError)
	if ok {
		ce.Details = append([]string{ce.Message}, ce.Details...)
		ce.Message = msg
		return ce
	}
	return &CError{
		Code:    http.StatusInternalServerError,
		Message: msg,
		Details: []string{err.Error()},
		stack:   util.Callers(),
	}
}

// IsRecordNotFound - check if an error is a RecordNotFound error
func IsRecordNotFound(err error) bool {
	return err.Error() == RecordNotFound
}
