package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
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

// RestoreLibrary deserializes a Library from a *bufio.Reader
func RestoreLibrary(reader *bufio.Reader) (*Library, Error) {
	library := NewLibrary()

	numRecords, err := ReadInt(reader)

	if err != nil || numRecords < 0 {
		return nil, NewlineError(ErrInvalidFile)
	}

	maxID := 0

	for i := 0; i < numRecords; i++ {
		record, err := RestoreRecord(reader)

		if err != nil {
			return nil, err
		}

		// no duplicate records allowed
		if _, ok := library.byID[record.id]; ok {
			return nil, NewlineError(ErrInvalidFile)
		} else if _, ok := library.byTitle[record.title]; ok {
			return nil, NewlineError(ErrInvalidFile)
		}

		library.byTitle[record.title] = record
		library.byID[record.id] = record

		if record.id > maxID {
			maxID = record.id
		}
	}

	library.nextID = maxID + 1

	return library, nil
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

// Save serializes a Library to an io.Writer in a format suitable for recovery
func (l *Library) Save(writer io.Writer) {
	FprintfOrPanic(writer, "%d\n", len(l.byTitle))

	for _, record := range l.sortedRecords() {
		record.Save(writer)
	}
}

// FindString returns all a string of all Records whose title contains a given substring,
// case insensitively. The Records are sorted by title in ascending order.
func (l *Library) FindString(substr string) (string, Error) {
	re := regexp.MustCompile(fmt.Sprintf("(?i)%s", regexp.QuoteMeta(substr)))

	var matches []*Record

	for _, record := range l.byTitle {
		if re.MatchString(record.title) {
			matches = append(matches, record)
		}
	}

	if len(matches) == 0 {
		return "", NewlineError("No records contain that string!")
	}

	SortRecordsByTitle(matches)

	return SprintRecords(matches), nil
}

const msgLibraryEmpty = "Library is empty"

// ListRatings returns a string of all Records sorted by rating in descending order.
// Records with the same rating are sorted by title in ascending order.
func (l *Library) ListRatings() string {
	if len(l.byTitle) == 0 {
		return msgLibraryEmpty
	}

	records := l.sortedRecords()

	sort.SliceStable(records, func(i, j int) bool {
		return records[i].rating > records[j].rating
	})

	return SprintRecords(records)
}

func (l *Library) String() string {
	if len(l.byTitle) == 0 {
		return msgLibraryEmpty
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Library contains %d records:\n", len(l.byTitle)))
	builder.WriteString(SprintRecords(l.sortedRecords()))

	return builder.String()
}

func (l *Library) doClear() {
	l.byTitle = make(libraryByTitle)
	l.byID = make(libraryByID)
	l.nextID = 1
}

func (l *Library) sortedRecords() []*Record {
	recordSet := make([]*Record, 0, len(l.byTitle))

	for _, record := range l.byTitle {
		recordSet = append(recordSet, record)
	}

	SortRecordsByTitle(recordSet)

	return recordSet
}
