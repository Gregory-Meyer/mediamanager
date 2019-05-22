package main

import (
	"fmt"
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
