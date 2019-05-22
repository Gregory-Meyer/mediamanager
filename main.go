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

func printCatalog(_ *Library, cat *Catalog) Error {
	fmt.Println(cat)

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
