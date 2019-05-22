package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

var stdin *bufio.Reader

func start() int {
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
		// "sA": saveAll,
		// "rA": restoreAll,
	}

	Library := NewLibrary()
	cat := NewCatalog()
	stdin = bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\nEnter command: ")

		cmd, e := readCommand()

		if e != nil { // EOF, probably
			return 0
		}

		if cmd == "qq" {
			break
		}

		if command, ok := commands[cmd]; !ok {
			fmt.Println("Unrecognized command!")
			readLineOrPanic()
		} else if e := command(Library, cat); e != nil {
			if e.ShouldSkipNewline() {
				readLineOrPanic()
			}

			fmt.Println(e)
		}
	}

	fmt.Println("Done")

	return 0
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
	medium := readWordOrPanic()
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
	name := readWordOrPanic()
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

	newRating, err := readInt()

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
	name := readWordOrPanic()
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
	id, err := readInt()

	if err != nil {
		return nil, err
	}

	return library.FindRecordByID(id)
}

func readCollection(catalog *Catalog) (*Collection, Error) {
	name, err := readWord()

	// this shouldn't happen
	if err != nil {
		panic(err)
	}

	return catalog.FindCollection(name)
}

func readTitle() (string, Error) {
	line := readLineOrPanic()
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

func readLineOrPanic() string {
	line, err := readLine()

	if err != nil {
		panic(err)
	}

	return line
}

func readLine() (string, error) {
	line, err := stdin.ReadString('\n')

	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(line, "\n"), nil
}

func readCommand() (string, error) {
	const commandLength = 2

	var command strings.Builder

	for i := 0; i < commandLength; i++ {
		if _, e := skipWhitespace(); e != nil {
			return "", e
		}

		r, _, e := stdin.ReadRune()

		if e != nil {
			return "", e
		}

		command.WriteRune(r)
	}

	return command.String(), nil
}

func readWordOrPanic() string {
	word, err := readWord()

	if err != nil {
		panic(err)
	}

	return word
}

func readWord() (string, error) {
	if _, e := skipWhitespace(); e != nil {
		return "", e
	}

	var word strings.Builder

	for {
		r, _, err := stdin.ReadRune()

		if err != nil {
			return "", err
		} else if unicode.IsSpace(r) {
			err = stdin.UnreadRune()

			if err != nil {
				return "", err
			}

			return word.String(), nil
		}

		word.WriteRune(r)
	}
}

const errUnreadableInteger = "Could not read an integer!"

func readInt() (int, Error) {
	if _, err := skipWhitespace(); err != nil {
		return 0, NewlineError(errUnreadableInteger)
	}

	var idBuilder strings.Builder

	r, _, err := stdin.ReadRune()

	if err != nil || (r != '+' && r != '-' && !unicode.IsNumber(r)) {
		return 0, NewlineError(errUnreadableInteger)
	}

	for {
		idBuilder.WriteRune(r)
		r, _, err = stdin.ReadRune()

		if err != nil {
			return 0, NewlineError(errUnreadableInteger)
		} else if !unicode.IsNumber(r) {
			err = stdin.UnreadRune()

			if err != nil {
				return 0, NewlineError(errUnreadableInteger)
			}

			break
		}
	}

	idStr := idBuilder.String()
	id, e := strconv.Atoi(idStr)

	if e != nil {
		return 0, NewlineError(errUnreadableInteger)
	}

	return id, nil
}

func skipWhitespace() (int, error) {
	numRead := 0

	for {
		r, _, err := stdin.ReadRune()

		if err != nil {
			return numRead, err
		} else if !unicode.IsSpace(r) {
			err = stdin.UnreadRune()

			if err != nil {
				return 0, err
			}

			return numRead, nil
		}

		numRead++
	}
}

func main() {
	os.Exit(start())
}
