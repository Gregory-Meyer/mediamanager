package main

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"
)

type catalogCollections map[string]*Collection

// Catalog is a set of named Collections
type Catalog struct {
	collections catalogCollections
}

// NewCatalog creates a Catalog ready to track Collections
func NewCatalog() *Catalog {
	return &Catalog{make(catalogCollections)}
}

// RestoreCatalog deserializes a Catalog from a *bufio.Reader
func RestoreCatalog(reader *bufio.Reader, library *Library) (*Catalog, Error) {
	numCollections, err := ReadInt(reader)

	if err != nil || numCollections < 0 {
		return nil, NewlineError(ErrInvalidFile)
	}

	catalog := NewCatalog()

	for i := 0; i < numCollections; i++ {
		collection, err := RestoreCollection(reader, library)

		if err != nil {
			return nil, err
		}

		if _, ok := catalog.collections[collection.name]; ok {
			return nil, NewlineError(ErrInvalidFile)
		}

		catalog.collections[collection.name] = collection
	}

	return catalog, nil
}

const errNoSuchCollection = "No collection with that name!"

// FindCollection indexes into a Collection by its name
func (c *Catalog) FindCollection(name string) (*Collection, Error) {
	collection, ok := c.collections[name]

	if !ok {
		return nil, NewlineError(errNoSuchCollection)
	}

	return collection, nil
}

// AddCollection adds a Collection to a Catalog
func (c *Catalog) AddCollection(name string) Error {
	if _, ok := c.collections[name]; ok {
		return NewlineError("Catalog already has a collection with this name!")
	}

	c.collections[name] = NewCollection(name)

	return nil
}

// DeleteCollection removes a Collection from a Catalog
func (c *Catalog) DeleteCollection(name string) Error {
	collection, ok := c.collections[name]

	if !ok {
		return NewlineError(errNoSuchCollection)
	}

	clearCollection(collection)
	delete(c.collections, name)

	return nil
}

// Clear removes all Collections from a Catalog
func (c *Catalog) Clear() {
	for _, collection := range c.collections {
		clearCollection(collection)
	}

	c.collections = make(catalogCollections)
}

// Save serializes a Catalog to an io.Writer in a format suitable for recovery
func (c *Catalog) Save(writer io.Writer) {
	FprintfOrPanic(writer, "%d\n", len(c.collections))

	for _, collection := range c.sortedCollections() {
		collection.Save(writer)
	}
}

func (c *Catalog) String() string {
	if len(c.collections) == 0 {
		return "Catalog is empty"
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Catalog contains %d collections:", len(c.collections)))

	for _, collection := range c.sortedCollections() {
		builder.WriteRune('\n')
		builder.WriteString(collection.String())
	}

	return builder.String()
}

// Clear erases all Records from this Collection's set of members
func clearCollection(collection *Collection) {
	for _, record := range collection.members {
		record.numCollections--
	}

	collection.members = make(collectionMembers)
}

func (c *Catalog) sortedCollections() []*Collection {
	collectionSet := make([]*Collection, 0, len(c.collections))

	for _, collection := range c.collections {
		collectionSet = append(collectionSet, collection)
	}

	sort.Slice(collectionSet, func(i, j int) bool {
		return collectionSet[i].name < collectionSet[j].name
	})

	return collectionSet
}
