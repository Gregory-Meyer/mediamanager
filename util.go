// MIT License
//
// Copyright (c) 2019 Gregory Meyer
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

// ErrInvalidFile is the error message when an rA command encounters a malformed file.
const ErrInvalidFile = "Invalid data found in file!"

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

// ReadLine reads until the next newline character or EOF, discarding the suffix
// Panics if an error is encountered
func ReadLine(reader *bufio.Reader) string {
	line, err := reader.ReadString('\n')

	if err != nil {
		if err == io.EOF {
			return line
		}

		panic(err)
	}

	return strings.TrimSuffix(line, "\n")
}

// ReadWord skips whitespace, then reads until, but not including, the next whitespace character
// Panics if an error is encountered
func ReadWord(reader *bufio.Reader) string {
	SkipWhitespace(reader)

	var word strings.Builder

	for {
		r, _, err := reader.ReadRune()

		if err != nil {
			if err == io.EOF {
				break
			}

			panic(err)
		} else if unicode.IsSpace(r) {
			err = reader.UnreadRune()

			if err != nil {
				panic(err)
			}

			break
		}

		word.WriteRune(r)
	}

	built := word.String()

	if len(built) == 0 {
		panic("didn't read any runes before EOF")
	}

	return word.String()
}

const errUnreadableInteger = "Could not read an integer value!"

// ReadInt skips whitespace, then reads up until the next non-numeric character
// Panics if an error is encountered while extracting input
func ReadInt(reader *bufio.Reader) (int, Error) {
	SkipWhitespace(reader)

	var idBuilder strings.Builder

	r, _, err := reader.ReadRune()

	if err != nil {
		if err == io.EOF {
			return 0, NewlineError(errUnreadableInteger)
		}

		panic(err)
	}

	if r != '+' && r != '-' && !unicode.IsNumber(r) {
		return 0, NewlineError(errUnreadableInteger)
	}

	for {
		idBuilder.WriteRune(r)
		r, _, err = reader.ReadRune()

		if err != nil {
			if err == io.EOF {
				break
			}

			panic(err)
		} else if !unicode.IsNumber(r) {
			err = reader.UnreadRune()

			if err != nil {
				panic(err)
			}

			break
		}
	}

	idStr := idBuilder.String()
	id, e := strconv.Atoi(idStr)

	if e != nil {
		return 0, NewlineError(errUnreadableInteger)
	}

	return id, nil
}

// SkipWhitespace reads up until, but not including, the next non-whitespace character
// Panics if an error is encountered
func SkipWhitespace(reader *bufio.Reader) {
	for {
		r, _, err := reader.ReadRune()

		if err != nil {
			if err == io.EOF {
				return
			}

			panic(err)
		} else if !unicode.IsSpace(r) {
			err = reader.UnreadRune()

			if err != nil {
				panic(err)
			}

			return
		}
	}
}
