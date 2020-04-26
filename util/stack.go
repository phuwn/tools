package util

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/phuwn/musichub/pkg/util"
)

// Frame - contains the current stack trace frame's information
type Frame struct {
	PC       uintptr
	FileName string
	Line     int
	FuncName string
}

func newFrame(f uintptr) *Frame {
	frame := &Frame{
		PC:       f - 1,
		FileName: "unknown",
		Line:     0,
		FuncName: "unknown",
	}

	fn := runtime.FuncForPC(frame.PC)
	if fn == nil {
		return frame
	}

	frame.FuncName = fn.Name()
	frame.FileName, frame.Line = fn.FileLine(frame.PC)
	return frame
}

// Format - implementation of Formatter interface of "fmt" package to print out frame
func (f Frame) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			io.WriteString(s, f.FuncName)
			io.WriteString(s, "\n\t")
			io.WriteString(s, f.FileName)
		default:
			io.WriteString(s, util.CleanPath(f.FileName))
		}
	case 'd':
		io.WriteString(s, strconv.Itoa(f.Line))
	case 'n':
		io.WriteString(s, formatFuncName(f.FuncName))
	case 'v':
		f.Format(s, 's')
		io.WriteString(s, ":")
		f.Format(s, 'd')
		if !s.Flag('+') {
			io.WriteString(s, " ")
			f.Format(s, 'n')
		}
	}
}

// Stack represents a stack of program counters.
type Stack []uintptr

// Format - implementation of Formatter interface of "fmt" package for custom print
func (s *Stack) Format(st fmt.State, verb rune) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(st, "\n failed to get current directory "+err.Error())
		return
	}
	switch verb {
	case 'v':
		switch {
		case st.Flag('+'):
			for _, pc := range *s {
				frame := newFrame(pc)
				if !strings.Contains(frame.FileName, currentDir) {
					break
				}
				fmt.Fprintf(st, "\n%+v", frame)
			}
		}
	}
}

// Caller - get current frame info
func Caller() *Frame {
	frame := &Frame{
		PC:       0,
		FileName: "unknown",
		Line:     0,
		FuncName: "unknown",
	}
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return frame
	}
	frame.PC, frame.FileName, frame.Line = pc, file, line

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return frame
	}
	frame.FuncName = fn.Name()

	return frame
}

// Callers - get current stack trace
func Callers() *Stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st Stack = pcs[0:n]
	return &st
}
