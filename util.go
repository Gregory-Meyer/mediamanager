package main

import (
	"fmt"
	"io"
)

// Error combines error reporting with newline handling
type Error interface {
	error
	ShouldSkipNewline() bool
}

// NewlineError indicates a newline should be skipped if caught
type NewlineError string

func (err NewlineError) Error() string {
	return string(err)
}

// ShouldSkipNewline returns true
func (NewlineError) ShouldSkipNewline() bool {
	return true
}

// RegularError indicates that a newline should not be skipped if caught
type RegularError string

func (err RegularError) Error() string {
	return string(err)
}

// ShouldSkipNewline returns false
func (RegularError) ShouldSkipNewline() bool {
	return false
}

// FprintfOrPanic panics if fmt.Fprintf returns an error
func FprintfOrPanic(writer io.Writer, format string, args ...interface{}) {
	_, err := fmt.Fprintf(writer, format, args...)

	if err != nil {
		panic(err)
	}
}
