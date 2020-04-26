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

type cError struct {
	code    int         `json:"code"`
	message string      `json:"message"`
	details []string    `json:"details"`
	stack   *util.Stack `json:"-"`
}

func (e cError) Error() string {
	return e.message
}

func (e cError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, e.message)
			if len(e.details) > 0 {
				io.WriteString(s, "\nDetails:\n- "+strings.Join(e.details, "\n- "))
			}
			e.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.message)
	case 'q':
		fmt.Fprintf(s, "%q", e.message)
	}
}

// New - generate a new custom error with input info and internal server status as default status
func New(msg string, args ...interface{}) error {
	return &cError{
		code:    http.StatusInternalServerError,
		message: fmt.Sprintf(msg, args...),
		stack:   util.Callers(),
	}
}

// NewResp - generate a new custom error with provided info
func NewResp(code int, msg string, details ...string) error {
	return &cError{
		code:    code,
		message: msg,
		details: details,
		stack:   util.Callers(),
	}
}

// AddDetails - add details to an error
func AddDetails(err error, details ...string) error {
	if err == nil {
		return nil
	}
	ce, ok := err.(*cError)
	if ok {
		ce.details = append(ce.details, details...)
		return ce
	}

	return &cError{
		code:    http.StatusInternalServerError,
		message: err.Error(),
		details: details,
		stack:   util.Callers(),
	}
}

// UpdateCode - update status of an error
func UpdateCode(err error, code int) error {
	if err == nil {
		return nil
	}
	ce, ok := err.(*cError)
	if ok {
		ce.code = code
		return ce
	}

	return &cError{
		code:    code,
		message: err.Error(),
		stack:   util.Callers(),
	}
}

// IsRecordNotFound - check if an error is a RecordNotFound error
func IsRecordNotFound(err error) bool {
	return err.Error() == RecordNotFound
}
