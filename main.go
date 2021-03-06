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
	"os"
	"strings"
)

var stdin *bufio.Reader

func main() {
	commands := map[string]func(*Library, *Catalog) Error{
		"fr": findRecord,
		"pr": printRecord,
		"pc": printCollection,
		"pL": printLibrary,
		"pC": printCatalog,
		"pa": printAllocations,
		"ar": addRecord,
		"ac": addCollection,
		"am": addMember,
		"mr": modifyRating,
		"dr": deleteRecord,
		"dc": deleteCollection,
		"dm": deleteMember,
		"cL": clearLibrary,
		"cC": clearCatalog,
		"cA": clearAll,
		"sA": saveAll,
		"rA": restoreAll,
		"fs": findString,
		"lr": listRatings,
		"cs": collectionStatistics,
		"cc": combineCollections,
		"mt": modifyTitle,
	}

	library := NewLibrary()
	catalog := NewCatalog()
	stdin = bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\nEnter command: ")

		cmd := readCommand()

		if cmd == "qq" {
			break
		}

		if command, ok := commands[cmd]; !ok {
			fmt.Println("Unrecognized command!")
			ReadLine(stdin)
		} else if e := command(library, catalog); e != nil {
			if e.ShouldSkipNewline() {
				ReadLine(stdin)
			}

			fmt.Println(e)
		}
	}

	_ = clearAll(library, catalog)
	fmt.Println("Done")
}

func findRecord(library *Library, _ *Catalog) Error {
	record, err := readRecordByTitle(library)

	if err != nil {
		return err
	}

	fmt.Println(record)

	return nil
}

func printRecord(library *Library, _ *Catalog) Error {
	record, err := readRecordByID(library)

	if err != nil {
		return err
	}

	fmt.Println(record)

	return nil
}

func printCollection(_ *Library, catalog *Catalog) Error {
	collection, err := readCollection(catalog)

	if err != nil {
		return err
	}

	fmt.Println(collection)

	return nil
}

func printLibrary(library *Library, _ *Catalog) Error {
	fmt.Println(library)

	return nil
}

func printCatalog(_ *Library, catalog *Catalog) Error {
	fmt.Println(catalog)

	return nil
}

func printAllocations(library *Library, catalog *Catalog) Error {
	fmtStr := `Memory allocations:
Records: %d
Collections: %d
`
	fmt.Printf(fmtStr, library.NumRecords(), catalog.NumCollections())

	return nil
}

func addRecord(library *Library, _ *Catalog) Error {
	medium := ReadWord(stdin)
	title, err := readTitle()

	if err != nil {
		return err
	}

	id, err := library.AddRecord(medium, title)

	if err != nil {
		return err
	}

	fmt.Printf("Record %d added\n", id)

	return nil
}

func addCollection(_ *Library, catalog *Catalog) Error {
	name := ReadWord(stdin)
	err := catalog.AddCollection(name)

	if err != nil {
		return err
	}

	fmt.Printf("Collection %s added\n", name)

	return nil
}

func addMember(library *Library, catalog *Catalog) Error {
	collection, err := readCollection(catalog)

	if err != nil {
		return err
	}

	record, err := readRecordByID(library)

	if err != nil {
		return err
	}

	err = collection.AddMember(record)

	if err != nil {
		return err
	}

	fmt.Printf("Member %d %s added\n", record.ID(), record.Title())

	return nil
}

func modifyRating(library *Library, _ *Catalog) Error {
	record, err := readRecordByID(library)

	if err != nil {
		return err
	}

	newRating, err := ReadInt(stdin)

	if err != nil {
		return err
	}

	err = record.SetRating(newRating)

	if err != nil {
		return err
	}

	fmt.Printf("Rating for record %d changed to %d\n", record.ID(), newRating)

	return nil
}

func deleteRecord(library *Library, _ *Catalog) Error {
	title, err := readTitle()

	if err != nil {
		return err
	}

	record, err := library.DeleteRecord(title)

	if err != nil {
		return err
	}

	fmt.Printf("Record %d %s deleted\n", record.ID(), record.Title())

	return nil
}

func deleteCollection(_ *Library, catalog *Catalog) Error {
	name := ReadWord(stdin)
	err := catalog.DeleteCollection(name)

	if err != nil {
		return err
	}

	fmt.Printf("Collection %s deleted\n", name)

	return nil
}

func deleteMember(library *Library, catalog *Catalog) Error {
	collection, err := readCollection(catalog)

	if err != nil {
		return err
	}

	record, err := readRecordByID(library)

	if err != nil {
		return err
	}

	err = collection.DeleteMember(record)

	if err != nil {
		return err
	}

	fmt.Printf("Member %d %s deleted\n", record.ID(), record.Title())

	return nil
}

func clearLibrary(library *Library, catalog *Catalog) Error {
	err := library.Clear(catalog)

	if err != nil {
		return err
	}

	fmt.Println("All records deleted")

	return nil
}

func clearCatalog(_ *Library, catalog *Catalog) Error {
	catalog.Clear()
	fmt.Println("All collections deleted")

	return nil
}

func clearAll(library *Library, catalog *Catalog) Error {
	library.ClearAll(catalog)
	fmt.Println("All data deleted")

	return nil
}

const errUnopenableFile = "Could not open file!"

func saveAll(library *Library, catalog *Catalog) Error {
	filename := ReadWord(stdin)
	file, err := os.Create(filename)

	if err != nil {
		return NewlineError(errUnopenableFile)
	}

	defer file.Close()

	library.Save(file)
	catalog.Save(file)

	fmt.Println("Data saved")

	return nil
}

func restoreAll(library *Library, catalog *Catalog) Error {
	filename := ReadWord(stdin)
	file, err := os.Open(filename)

	if err != nil {
		return NewlineError(errUnopenableFile)
	}

	defer file.Close()
	reader := bufio.NewReader(file)

	newLibrary, parseErr := RestoreLibrary(reader)

	if parseErr != nil {
		return parseErr
	}

	newCatalog, parseErr := RestoreCatalog(reader, newLibrary)

	if parseErr != nil {
		return parseErr
	}

	*library = *newLibrary
	*catalog = *newCatalog

	fmt.Println("Data loaded")

	return nil
}

func findString(library *Library, _ *Catalog) Error {
	substr := ReadWord(stdin)
	matches, err := library.FindString(substr)

	if err != nil {
		return err
	}

	fmt.Println(matches)

	return nil
}

func listRatings(library *Library, _ *Catalog) Error {
	fmt.Println(library.ListRatings())

	return nil
}

func collectionStatistics(library *Library, catalog *Catalog) Error {
	numOne, numMany, total := catalog.CollectionStatistics()
	numRecords := library.NumRecords()

	// could use string concatenation instead here
	fmtStr := `%d out of %d Records appear in at least one Collection
%d out of %d Records appear in more than one Collection
Collections contain a total of %d Records
`
	fmt.Printf(fmtStr, numOne, numRecords, numMany, numRecords, total)

	return nil
}

func combineCollections(_ *Library, catalog *Catalog) Error {
	firstSrc, err := readCollection(catalog)

	if err != nil {
		return err
	}

	secondSrc, err := readCollection(catalog)

	if err != nil {
		return err
	}

	dstName := ReadWord(stdin)

	err = catalog.CombineCollections(firstSrc, secondSrc, dstName)

	if err != nil {
		return err
	}

	fmt.Printf("Collections %s and %s combined into new collection %s\n",
		firstSrc.Name(), secondSrc.Name(), dstName)

	return nil
}

func modifyTitle(library *Library, _ *Catalog) Error {
	record, err := readRecordByID(library)

	if err != nil {
		return err
	}

	newTitle, err := readTitle()

	if err != nil {
		return err
	}

	err = library.ModifyTitle(record, newTitle)

	if err != nil {
		return err
	}

	fmt.Printf("Title for record %d changed to %s\n", record.ID(), newTitle)

	return nil
}

func readRecordByTitle(library *Library) (*Record, Error) {
	title, err := readTitle()

	if err != nil {
		return nil, err
	}

	return library.FindRecordByTitle(title)
}

func readRecordByID(library *Library) (*Record, Error) {
	id, err := ReadInt(stdin)

	if err != nil {
		return nil, err
	}

	return library.FindRecordByID(id)
}

func readCollection(catalog *Catalog) (*Collection, Error) {
	name := ReadWord(stdin)

	return catalog.FindCollection(name)
}

func readTitle() (string, Error) {
	line := ReadLine(stdin)
	fields := strings.Fields(line)

	if len(fields) == 0 {
		return "", RegularError("Could not read a title!")
	}

	var title strings.Builder

	title.WriteString(fields[0])

	for _, word := range fields[1:] {
		title.WriteRune(' ')
		title.WriteString(word)
	}

	return title.String(), nil
}

func readCommand() string {
	const commandLength = 2

	var command strings.Builder

	for i := 0; i < commandLength; i++ {
		SkipWhitespace(stdin)
		r, _, err := stdin.ReadRune()

		if err != nil {
			panic(err)
		}

		command.WriteRune(r)
	}

	return command.String()
}
