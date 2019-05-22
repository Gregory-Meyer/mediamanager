package main

import (
	"fmt"
	"io"
)

// FprintfOrPanic panics if fmt.Fprintf returns an error
func FprintfOrPanic(writer io.Writer, format string, args ...interface{}) {
	_, err := fmt.Fprintf(writer, format, args...)

	if err != nil {
		panic(err)
	}
}
