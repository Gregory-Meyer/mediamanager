package main

import (
	"fmt"
	"sort"
	"strings"
)

type libraryByTitle map[string]*Record
type libraryByID map[int]*Record

// Library is a set of Records that can be indexed by title or by ID
type Library struct {
	byTitle libraryByTitle
	byID    libraryByID
	nextID  int
}

// NewLibrary creates an empty Library ready to track Records
func NewLibrary() *Library {
	return &Library{make(libraryByTitle), make(libraryByID), 1}
}

const errNoSuchRecordTitle = "No record with that title!"

// FindRecordByTitle indexes into a Library's set of Records by title
func (l *Library) FindRecordByTitle(title string) (*Record, Error) {
	record, ok := l.byTitle[title]

	if !ok {
		return nil, RegularError(errNoSuchRecordTitle)
	}

	return record, nil
}

// FindRecordByID indexes into a Library's set of Records by ID
func (l *Library) FindRecordByID(id int) (*Record, Error) {
	record, ok := l.byID[id]

	if !ok {
		return nil, NewlineError("No record with that ID!")
	}

	return record, nil
}

// AddRecord adds a Record into the Library
func (l *Library) AddRecord(medium, title string) (int, Error) {
	if _, ok := l.byTitle[title]; ok {
		return 0, RegularError("Library already has a record with this title!")
	}

	id := l.nextID
	record := NewRecord(medium, title, id)
	l.nextID++

	l.byTitle[title] = record
	l.byID[id] = record

	return id, nil
}

// DeleteRecord erases a Record from this Library's set
func (l *Library) DeleteRecord(title string) (*Record, Error) {
	record, ok := l.byTitle[title]

	if !ok {
		return nil, RegularError(errNoSuchRecordTitle)
	}

	if record.numCollections > 0 {
		return nil, RegularError("Cannot delete a record that is a member of a collection!")
	}

	delete(l.byTitle, record.title)
	delete(l.byID, record.id)

	return record, nil
}

// Clear erases all Records from this Library's set
func (l *Library) Clear(catalog *Catalog) Error {
	for _, collection := range catalog.collections {
		if len(collection.members) > 0 {
			return NewlineError("Cannot clear all records unless all collections are empty!")
		}
	}

	l.doClear()

	return nil
}

// ClearAll erases all Records from a Library and all Collections from a Catalog
func (l *Library) ClearAll(catalog *Catalog) {
	catalog.Clear()
	l.doClear()
}

func (l *Library) doClear() {
	l.byTitle = make(libraryByTitle)
	l.byID = make(libraryByID)
	l.nextID = 1
}

func (l *Library) String() string {
	if len(l.byTitle) == 0 {
		return "Library is empty"
	}

	titleSet := make([]string, 0, len(l.byTitle))

	for title := range l.byTitle {
		titleSet = append(titleSet, title)
	}

	sort.Strings(titleSet)

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Library contains %d records:", len(titleSet)))

	for _, title := range titleSet {
		builder.WriteRune('\n')
		builder.WriteString(l.byTitle[title].String())
	}

	return builder.String()
}
