// Package status provides basic facility to print status line to the terminal.
//
// Typical usage:
//
//  line := new(status.Line)
//  defer line.Done()
//  for i := 0; i < total; i++ {
//      line.Printf("step %d out of %d", i, total)
//      ...
//  }
//
// This package outputs VT100 escape sequences, so its output may break if your
// terminal does not support them.
package status

import (
	"fmt"
	"os"
	"sync"

	"golang.org/x/term"
)

// Line can be used to print a single status line, overwriting it on each
// Print/Printf call. Finish with Done call to move cursor to the new line.
//
// By default Line writes to os.Stdout, use SetOutput method to override.
//
// If Line's output is not connected to a terminal (for example, if program's
// output is redirected to a file), Print/Printf calls do nothing.
type Line struct {
	once sync.Once
	w    *os.File
	noop bool
	buf  []byte
}

// SetOutput overrides the default os.Stdout output. Can only be called before
// any other methods.
func (l *Line) SetOutput(f *os.File) {
	if l.w != nil {
		panic("Line.SetOutput must be called before any of other Line methods")
	}
	l.w = f
	l.init()
}

func (l *Line) init() {
	l.once.Do(func() {
		if l.w == nil {
			l.w = os.Stdout
		}
		l.noop = !term.IsTerminal(int(l.w.Fd()))
	})
}

// Print works like fmt.Print. Because Line is expected to write over a single
// line, including any newlines will break output.
func (l *Line) Print(a ...interface{}) (n int, err error) {
	l.init()
	if l.noop {
		return 0, nil
	}
	l.buf = append(l.buf[:0], escPrefix...)
	l.buf = append(l.buf, fmt.Sprint(a...)...)
	return l.w.Write(l.buf)
}

// Printf works like fmt.Printf. Because Line is expected to write over a
// single line, including any newlines will break output.
func (l *Line) Printf(format string, a ...interface{}) (n int, err error) {
	l.init()
	if l.noop {
		return 0, nil
	}
	l.buf = append(l.buf[:0], escPrefix...)
	l.buf = append(l.buf, fmt.Sprintf(format, a...)...)
	return l.w.Write(l.buf)
}

// Done writes a single newline.
func (l *Line) Done() {
	l.init()
	if l.noop {
		return
	}
	l.w.WriteString("\n")
}

const escPrefix = "\x1b[2K\r" // "\x1b[2K\x1b[G"
