package handler

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/phuwn/tools/errors"
	"github.com/phuwn/tools/log"
)

// JSONResponse - json response framework
type JSONResponse struct {
	Code  int            `json:"code"`
	Error *errors.CError `json:"error,omitempty"`
	Data  interface{}    `json:"data,omitempty"`
}

func errFormat(err error, c echo.Context) *errors.CError {
	ce, ok := err.(*errors.CError)
	if ok {
		return ce
	}
	switch err {
	case echo.ErrNotFound:
		return &errors.CError{Code: 404, Message: fmt.Sprintf(`url not found: %s`, c.Path())}
	default:
		return &errors.CError{Code: 500, Message: err.Error()}
	}
}

// JSON - json response handler
func JSON(c echo.Context, code int, data interface{}) error {
	return c.JSON(code, &JSONResponse{code, nil, data})
}

// JSONError - json error response handler
func JSONError(err error, c echo.Context) {
	er := errFormat(err, c)
	if er.Code == 500 {
		log.Error(err)
	}

	c.JSON(er.Code, &JSONResponse{er.Code, er, nil})
}
