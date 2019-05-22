package main

import (
	"bufio"
	"fmt"
	"io"
)

// Record is a piece of media
type Record struct {
	medium         string
	title          string
	rating         int
	id             int
	numCollections int
}

// NewRecord creates a Record
func NewRecord(medium, title string, id int) *Record {
	return &Record{medium, title, 0, id, 0}
}

// RestoreRecord deserializes a Record from a *bufio.Reader
func RestoreRecord(reader *bufio.Reader) (*Record, Error) {
	id, err := ReadInt(reader)

	if err != nil || id < 1 {
		return nil, NewlineError(ErrInvalidFile)
	}

	medium := ReadWord(reader)

	// EOF
	if len(medium) == 0 {
		return nil, NewlineError(ErrInvalidFile)
	}

	rating, err := ReadInt(reader)

	if err != nil || rating < 0 || rating > 5 {
		return nil, NewlineError(ErrInvalidFile)
	}

	SkipWhitespace(reader)
	title := ReadLine(reader)

	if len(title) == 0 {
		return nil, NewlineError(ErrInvalidFile)
	}

	return &Record{medium, title, rating, id, 0}, nil
}

// ID gives the ID of this Record, which starts at 1 and goes up from there
func (r *Record) ID() int {
	return r.id
}

// Title gives the title of this Record, which is a string with at most one
// space between words
func (r *Record) Title() string {
	return r.title
}

// SetRating sets the rating of this Record
// Ratings are between 1 and 5, inclusive
func (r *Record) SetRating(newRating int) Error {
	if newRating < 1 || newRating > 5 {
		return NewlineError("Rating is out of range!")
	}

	r.rating = newRating

	return nil
}

// Save serializes a Record to an io.Writer in a format suitable for recovery
func (r *Record) Save(writer io.Writer) {
	FprintfOrPanic(writer, "%d %s %d %s\n", r.id, r.medium, r.rating, r.title)
}

func (r *Record) String() string {
	if r.rating == 0 {
		return fmt.Sprintf("%d: %s u %s", r.id, r.medium, r.title)
	}

	return fmt.Sprintf("%d: %s %d %s", r.id, r.medium, r.rating, r.title)
}
