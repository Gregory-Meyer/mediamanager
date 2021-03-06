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

const errDuplicateCollection = "Catalog already has a collection with this name!"

// NumCollections returns the number of Collections in this Catalog
func (c *Catalog) NumCollections() int {
	return len(c.collections)
}

// AddCollection adds a Collection to a Catalog
func (c *Catalog) AddCollection(name string) Error {
	if _, ok := c.collections[name]; ok {
		return NewlineError(errDuplicateCollection)
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

	*c = *NewCatalog()
}

// Save serializes a Catalog to an io.Writer in a format suitable for recovery
func (c *Catalog) Save(writer io.Writer) {
	FprintfOrPanic(writer, "%d\n", len(c.collections))

	for _, collection := range c.sortedCollections() {
		collection.Save(writer)
	}
}

// CollectionStatistics computes the number of Records that appear in at least
// one Collection, the number of Records that appear in more than one
// Collection, and the total of Records contained by Collections
func (c *Catalog) CollectionStatistics() (numOne int, numMany int, total int) {
	counts := make(map[int]int) // record ID -> number of collections containing it

	// fabled double for loop for minimum performance
	for _, collection := range c.collections {
		total += len(collection.members)

		for _, record := range collection.members {
			// no way to avoid double lookup here, but Go makes map lookups cheap
			if prevCount, ok := counts[record.id]; ok {
				counts[record.id] = prevCount + 1

				if prevCount == 1 {
					numMany++
				}
			} else {
				counts[record.id] = 1
				numOne++
			}
		}
	}

	return numOne, numMany, total
}

// CombineCollections combines two source Collections into a destination
// Collection with a new name, leaving the two source Collections unmodified
func (c *Catalog) CombineCollections(firstSrc, secondSrc *Collection, dstName string) Error {
	if _, ok := c.collections[dstName]; ok {
		return NewlineError(errDuplicateCollection)
	}

	dst := NewCollection(dstName)
	c.collections[dstName] = dst

	for _, record := range firstSrc.members {
		_ = dst.AddMember(record)
	}

	for _, record := range secondSrc.members {
		_ = dst.AddMember(record)
	}

	return nil
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
