package main

import (
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
