package main

import (
	"fmt"
	"io"
	"sort"
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

	for _, title := range c.sortedMemberTitles() {
		FprintfOrPanic(writer, "%s\n", title)
	}
}

func (c *Collection) String() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Collection %s contains:", c.name))

	if len(c.members) == 0 {
		builder.WriteString(" None")
	} else {
		for _, r := range c.members {
			builder.WriteRune('\n')
			builder.WriteString(r.String())
		}
	}

	return builder.String()
}

func (c *Collection) sortedMemberTitles() []string {
	titleSet := make([]string, 0, len(c.members))

	for _, record := range c.members {
		titleSet = append(titleSet, record.title)
	}

	sort.Strings(titleSet)

	return titleSet
}
