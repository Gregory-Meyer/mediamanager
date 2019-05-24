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
	"strings"
)

type collectionMembers map[int]*Record

// Collection is a named set of Records
type Collection struct {
	name    string
	members collectionMembers
}

// NewCollection creates a Collection
func NewCollection(name string) *Collection {
	return &Collection{name, make(collectionMembers)}
}

// RestoreCollection deserializes a Collection from a *bufio.Reader
func RestoreCollection(reader *bufio.Reader, library *Library) (*Collection, Error) {
	name := ReadWord(reader)

	// EOF
	if len(name) == 0 {
		return nil, NewlineError(ErrInvalidFile)
	}

	collection := NewCollection(name)
	numMembers, err := ReadInt(reader)

	if err != nil || numMembers < 0 {
		return nil, NewlineError(ErrInvalidFile)
	}

	ReadLine(reader)

	for i := 0; i < numMembers; i++ {
		title := ReadLine(reader)
		record, ok := library.byTitle[title]

		if !ok {
			return nil, NewlineError(ErrInvalidFile)
		}

		if _, ok := collection.members[record.id]; ok {
			return nil, NewlineError(ErrInvalidFile)
		}

		collection.members[record.id] = record
		record.numCollections++
	}

	return collection, nil
}

// AddMember inserts a Record into this Collection's set of members
func (c *Collection) AddMember(record *Record) Error {
	if _, ok := c.members[record.id]; ok {
		return NewlineError("Record is already a member in the collection!")
	}

	c.members[record.id] = record
	record.numCollections++

	return nil
}

// DeleteMember erases a Record from this Collection's set of members
func (c *Collection) DeleteMember(record *Record) Error {
	if _, ok := c.members[record.id]; !ok {
		return NewlineError("Record is not a member in the collection!")
	}

	delete(c.members, record.id)
	record.numCollections--

	return nil
}

// Save serializes a Collection to an io.Writer in a format suitable for recovery
func (c *Collection) Save(writer io.Writer) {
	FprintfOrPanic(writer, "%s %d\n", c.name, len(c.members))

	for _, record := range c.sortedMembers() {
		FprintfOrPanic(writer, "%s\n", record.title)
	}
}

// Name returns the name of this Collection
func (c *Collection) Name() string {
	return c.name
}

func (c *Collection) String() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Collection %s contains:", c.name))

	if len(c.members) == 0 {
		builder.WriteString(" None")
	} else {
		builder.WriteRune('\n')
		builder.WriteString(SprintRecords(c.sortedMembers()))
	}

	return builder.String()
}

func (c *Collection) sortedMembers() []*Record {
	memberSet := make([]*Record, 0, len(c.members))

	for _, record := range c.members {
		memberSet = append(memberSet, record)
	}

	SortRecordsByTitle(memberSet)

	return memberSet
}
