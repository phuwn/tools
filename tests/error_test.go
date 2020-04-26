package tests

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/phuwn/tools/errors"
)

func Test_cError_Print(t *testing.T) {
	cErr := errors.New("testing")
	currentDir, err := os.Getwd()
	if err != nil {
		t.Errorf("failed to get current directory %s\n", err.Error())
		return
	}
	tests := []struct {
		name   string
		format string
		want   string
	}{
		{
			name:   "message only",
			format: "%v",
			want:   "testing",
		},
		{
			name:   "log full error with tracing",
			format: "%+v",
			want:   fmt.Sprintf("testing\ngithub.com/phuwn/tools/tests.Test_cError_Print\n\t%s/error_test.go:13", currentDir),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			_, err := fmt.Fprintf(&b, tt.format, cErr)
			if err != nil {
				t.Errorf("print process got error %s", err.Error())
				return
			}
			if b.String() != tt.want {
				t.Errorf("strange output\nwant: \n%s\ngot: \n%s\n", tt.want, b.String())
			}
		})
	}
}
