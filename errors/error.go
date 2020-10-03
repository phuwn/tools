package errors

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"regexp"

	"github.com/phuwn/tools/util"
)

var (
	re = regexp.MustCompile(`["` + "\t" + "\n]")
)

type MessageStack []string

func (ms MessageStack) Print() string {
	var res string
	for i := len(ms) - 1; i >= 0; i-- {
		res += "- " + ms[i] + "\n"
	}
	return res
}

func (ms MessageStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v', 's':
		for i := len(ms) - 1; i >= 0; i-- {
			io.WriteString(s, "\n- "+ms[i])
		}
	}
}

func (ms MessageStack) MarshalJSON() ([]byte, error) {
	res := "["
	for i := len(ms) - 1; i >= 0; i-- {
		res += `"` + re.ReplaceAllString(ms[i], "") + `"`
		if i != 0 {
			res += ","
		}
	}
	return []byte(res + "]"), nil
}

var (
	// RecordNotFound error message
	RecordNotFound = "record not found"

	InternalServerErrorJSON = `{"error":"internal server error"}`
)

// CError - Custom error for debugging and json responsing
type CError struct {
	Code      int          `json:"-"`
	Message   string       `json:"error"`
	Details   MessageStack `json:"details,omitempty"`
	callStack *util.CallStack
}

// Error - Return error content
func (e CError) Error() string {
	return e.Message
}

// Format - implement Formatter interface for custom print
func (e CError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, e.Message)
			if len(e.Details) > 0 {
				io.WriteString(s, "\nDetails:")
				e.Details.Format(s, 'v')
			}
			e.callStack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.Message)
	case 'q':
		fmt.Fprintf(s, "%q", e.Message)
	}
}

// MarshalJSON - implement Mashaler interface for response in JSON
func (e CError) MarshalJSON() ([]byte, error) {
	if e.Code == 500 {
		return []byte(InternalServerErrorJSON), nil
	}

	res := fmt.Sprintf(`{"error":"%s"`, e.Message)
	if len(e.Details) > 0 {
		details, _ := e.Details.MarshalJSON()
		res += fmt.Sprintf(`,"details":%s`, string(details))
	}
	return []byte(res + "}"), nil
}

// New - generate a new custom error with error message and status code, default code will be internal server status code
func New(msg string, code ...int) error {
	stt := http.StatusInternalServerError
	if len(code) > 0 {
		stt = code[0]
	}
	return &CError{
		Code:      stt,
		Message:   msg,
		callStack: util.Callers(),
	}
}

// Customize - customize normal error to CError, push new message ahead of the details queue
func Customize(err error, code int, msg string) error {
	if reflect.ValueOf(err).IsNil() {
		return nil
	}

	ce, ok := err.(*CError)
	if !ok {
		ce = &CError{
			Message:   err.Error(),
			callStack: util.Callers(),
		}
	}

	ce.Code = code
	if msg != "" {
		ce.Details = append(ce.Details, ce.Message)
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
