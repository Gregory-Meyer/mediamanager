package main

import (
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

const errNoSuchCollection = "Could not read an integer!"

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

	for _, name := range c.sortedCollectionNames() {
		c.collections[name].Save(writer)
	}
}

func (c *Catalog) String() string {
	if len(c.collections) == 0 {
		return "Catalog is empty"
	}

	nameSet := make([]string, 0, len(c.collections))

	for title := range c.collections {
		nameSet = append(nameSet, title)
	}

	sort.Strings(nameSet)

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Catalog contains %d collections:", len(c.collections)))

	for _, name := range nameSet {
		builder.WriteRune('\n')
		builder.WriteString(c.collections[name].String())
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

func (c *Catalog) sortedCollectionNames() []string {
	nameSet := make([]string, 0, len(c.collections))

	for name := range c.collections {
		nameSet = append(nameSet, name)
	}

	sort.Strings(nameSet)

	return nameSet
}
