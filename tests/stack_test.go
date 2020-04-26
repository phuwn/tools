package tests

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/phuwn/tools/util"
)

func Test_Frame_Print(t *testing.T) {
	frame := util.Caller()
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
			name:   "file clean path",
			format: "%s",
			want:   "tests/stack_test.go",
		},
		{
			name:   "file path detail",
			format: "%+s",
			want:   fmt.Sprintf("github.com/phuwn/tools/tests.Test_Frame_Print\n\t%s/stack_test.go", currentDir),
		},
		{
			name:   "code line only",
			format: "%d",
			want:   "13",
		},
		{
			name:   "func name only",
			format: "%n",
			want:   "Test_Frame_Print",
		},
		{
			name:   "short frame info",
			format: "%v",
			want:   "tests/stack_test.go:13 Test_Frame_Print",
		},
		{
			name:   "full frame info",
			format: "%+v",
			want:   fmt.Sprintf("github.com/phuwn/tools/tests.Test_Frame_Print\n\t%s/stack_test.go:13", currentDir),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			_, err := fmt.Fprintf(&b, tt.format, frame)
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
