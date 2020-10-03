package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func Test_CError_Print(t *testing.T) {
	cErr := New("testing")
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
			want:   fmt.Sprintf("testing\ngithub.com/phuwn/tools/errors.Test_CError_Print\n\t%s/error_test.go:13", currentDir),
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

func Test_CError_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		content error
		want    string
	}{
		{
			name:    "500 status code CError",
			content: New("simple error"),
			want:    InternalServerErrorJSON,
		},
		{
			name:    "simple CError",
			content: Customize(fmt.Errorf("token is invalid"), 401, "Unauthorized"),
			want:    `{"error":"Unauthorized","details":["token is invalid"]}`,
		},
		{
			name:    "detail is a JSON",
			content: Customize(fmt.Errorf(`{"error":"token is invalid"}`), 401, "Unauthorized"),
			want:    `{"error":"Unauthorized","details":["{error:token is invalid}"]}`,
		},
		{
			name:    "detail got endline",
			content: Customize(fmt.Errorf("token is \ninvalid"), 401, "Unauthorized"),
			want:    `{"error":"Unauthorized","details":["token is invalid"]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(tt.content)
			if err != nil {
				t.Errorf("unknown error occurs %v", err.Error())
				return
			}
			if string(b) != tt.want {
				t.Errorf("strange output\nwant: \n%s\ngot: \n%s\n", tt.want, string(b))
			}
		})
	}
}

func Test_CError_Customize(t *testing.T) {
	type args struct {
		err  error
		code int
		msg  string
	}
	tests := []struct {
		name string
		args args
		want *CError
	}{
		{
			"customize an error from system",
			args{fmt.Errorf("record not found"), 404, "failed to find user 1"},
			&CError{404, "failed to find user 1", MessageStack{"record not found"}, nil},
		},
		{
			"customize a CError",
			args{New("error msg 1"), 400, "error msg 2"},
			&CError{400, "error msg 2", MessageStack{"error msg 1"}, nil},
		},
		{
			"customize a CError with details",
			args{&CError{500, "error msg 1", MessageStack{"error detail 2", "error detail 1"}, nil}, 400, "error msg 2"},
			&CError{400,
				"error msg 2",
				MessageStack{
					"error detail 2",
					"error detail 1",
					"error msg 1",
				},
				nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Customize(tt.args.err, tt.args.code, tt.args.msg)
			if reflect.ValueOf(err).IsNil() != reflect.ValueOf(tt.want).IsNil() {
				t.Errorf("unknown error result\nwant: \n%v\ngot: \n%v\n", tt.want, err)
				return
			}
			if !reflect.ValueOf(tt.want).IsNil() {
				ce, ok := err.(*CError)
				if !ok {
					t.Errorf("unknown error occurs, cant convert error to CError\nwant: \n%v\ngot: \n%v\n", tt.want, err)
					return
				}
				if tt.want.Code != ce.Code {
					t.Errorf("wrong status code\nwant: \n%d\ngot: \n%d\n", tt.want.Code, ce.Code)
					return
				}
				if tt.want.Message != ce.Message {
					t.Errorf("unknown message\nwant: \n%s\ngot: \n%s\n", tt.want.Message, ce.Message)
					return
				}
				ad, bd := tt.want.Details.Print(), ce.Details.Print()
				if ad != bd {
					t.Errorf("unknown detail message\nwant: \n%s\ngot: \n%s\n", ad, bd)
				}
			}
		})
	}
}
