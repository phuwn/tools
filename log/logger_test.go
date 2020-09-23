package log

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/phuwn/tools/errors"
	"github.com/phuwn/tools/util"
)

func TestDefaultLogger_Error(t *testing.T) {
	util.MockRuntimeFunc()
	var (
		simpleError            = fmt.Errorf("simple error")
		customError            = errors.New("custom error")
		customErrorWithDetail  = errors.Customize(simpleError, 500, "new message")
		customErrorWithDetails = errors.Customize(errors.Customize(simpleError, 500, "new message 1"), 500, "new message 2")
	)
	tests := []struct {
		name   string
		locale *time.Location
		err    error
		args   []interface{}
		want   string
	}{
		{
			"print simple error only",
			nil,
			simpleError,
			nil,
			"2020/09/20 17:00:58UTC [ERROR] simple error\n",
		},
		{
			"print simple error with args",
			nil,
			simpleError,
			[]interface{}{"testing with", 4, "args", true},
			"2020/09/20 17:00:58UTC [ERROR] testing with 4 args true\nsimple error\n",
		},
		{
			"print custom error with only message",
			nil,
			customError,
			nil,
			"2020/09/20 17:00:58UTC [ERROR] custom error\ngithub.com/phuwn/tools/log.TestDefaultLogger_Error\n\t/Users/phuong/go/src/github.com/phuwn/tools/log/logger_test.go:17\n",
		},
		{
			"print custom error with 1 detail",
			nil,
			customErrorWithDetail,
			nil,
			"2020/09/20 17:00:58UTC [ERROR] new message\nDetails:\n- simple error\ngithub.com/phuwn/tools/log.TestDefaultLogger_Error\n\t/Users/phuong/go/src/github.com/phuwn/tools/log/logger_test.go:18\n",
		},
		{
			"print custom error with details and args",
			nil,
			customErrorWithDetails,
			[]interface{}{"testing with", 3.0, "args"},
			"2020/09/20 17:00:58UTC [ERROR] testing with 3 args\nnew message 2\nDetails:\n- new message 1\n- simple error\ngithub.com/phuwn/tools/log.TestDefaultLogger_Error\n\t/Users/phuong/go/src/github.com/phuwn/tools/log/logger_test.go:19\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			l := DefaultLogger{b, tt.locale}
			l.Error(tt.err, tt.args...)
			if b.String() != tt.want {
				t.Errorf("strange output\nwant: \n%v\ngot: \n%v\n", tt.want, b.String())
			}
		})
	}
}

func TestDefaultLogger_Info(t *testing.T) {
	util.MockRuntimeFunc()
	tests := []struct {
		name   string
		locale *time.Location
		args   []interface{}
		want   string
	}{
		{
			"print simple text",
			nil,
			[]interface{}{"simple text"},
			"2020/09/20 17:00:58UTC [INFO] simple text testing/testing.go:909 tRunner\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			l := DefaultLogger{b, tt.locale}
			l.Info(tt.args...)
			if b.String() != tt.want {
				t.Errorf("strange output\nwant: \n%v\ngot: \n%v\n", tt.want, b.String())
			}
		})
	}
}
