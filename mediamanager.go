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
	commands := map[string]func(*library, *catalog, *int) err{
		"fr": findRecord,
		"pr": printRecord,
		// "pc": printCollection,
		// "pL": printLibrary,
		// "pC": printCatalog,
		"ar": addRecord,
		// "ac": addCollection,
		// "am": addMember,
		// "mr": modifyRating,
		// "dr": deleteRecord,
		// "dc": deleteCollection,
		// "dm": deleteMember,
		// "cL": clearLibrary,
		// "cC": clearCatalog,
		// "cA": clearAll,
		// "sA": saveAll,
		// "rA": restoreAll,
	}

	lib := newLibrary()
	cat := newCatalog()
	nextID := 1
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
			readLine()
		} else if e := command(&lib, &cat, &nextID); e != nil {
			if e.ShouldSkipNewline() {
				readLine()
			}

			fmt.Println(e)
		}
	}

	fmt.Println("Done")

	return 0
}

type library struct {
	byTitle map[string]*record
	byID    map[int]*record
}

func newLibrary() library {
	return library{make(map[string]*record), make(map[int]*record)}
}

type catalog map[string]collection

func newCatalog() catalog {
	return make(catalog)
}

type collection map[string]*record

func newCollection() map[string]*record {
	return make(collection)
}

type record struct {
	medium string
	title  string
	rating int
	id     int
}

func newRecord(medium, title string, id int) record {
	return record{medium, title, 0, id}
}

func (r *record) String() string {
	if r.rating == 0 {
		return fmt.Sprintf("%d: %s u %s", r.id, r.medium, r.title)
	}

	return fmt.Sprintf("%d: %s %d %s", r.id, r.medium, r.rating, r.title)
}

type err interface {
	error
	ShouldSkipNewline() bool
}

type newlineErr struct {
	what string
}

func (e newlineErr) ShouldSkipNewline() bool {
	return true
}

func (e newlineErr) Error() string {
	return e.what
}

type regularErr struct {
	what string
}

func (e regularErr) ShouldSkipNewline() bool {
	return false
}

func (e regularErr) Error() string {
	return e.what
}

func findRecord(lib *library, _ *catalog, _ *int) err {
	title, e := readTitle()

	if e != nil {
		return e
	}

	if rec, ok := lib.byTitle[title]; ok {
		fmt.Println(rec)

		return nil
	}

	return regularErr{"No record with that title!"}
}

func printRecord(lib *library, _ *catalog, _ *int) err {
	id, e := readID()

	if e != nil {
		return e
	}

	if rec, ok := lib.byID[id]; ok {
		fmt.Println(rec)

		return nil
	}

	return newlineErr{"No record with that ID!"}
}

func addRecord(lib *library, _ *catalog, id *int) err {
	// should never fail unless EOF, which the spec says to ignore
	medium, _ := readWord()
	title, e := readTitle()

	if e != nil {
		return e
	}

	if _, ok := lib.byTitle[title]; ok {
		return regularErr{"Library already has a record with this title!"}
	}

	thisID := *id
	(*id)++

	rec := newRecord(medium, title, thisID)

	lib.byTitle[title] = &rec
	lib.byID[thisID] = &rec

	return nil
}

func readTitle() (string, err) {
	line, e := readLine()

	if e != nil {
		return "", newlineErr{e.Error()}
	}

	fields := strings.Fields(line)

	if len(fields) == 0 {
		return "", regularErr{"Could not read a title!"}
	}

	var title strings.Builder

	title.WriteString(fields[0])

	for _, word := range fields[1:] {
		title.WriteRune(' ')
		title.WriteString(word)
	}

	return title.String(), nil
}

func readLine() (string, error) {
	line, e := stdin.ReadString('\n')

	if e != nil {
		return "", e
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

func readWord() (string, error) {
	if _, e := skipWhitespace(); e != nil {
		return "", e
	}

	var word strings.Builder

	for {
		r, _, e := stdin.ReadRune()

		if e != nil {
			return "", e
		} else if unicode.IsSpace(r) {
			stdin.UnreadRune()

			return word.String(), nil
		}

		word.WriteRune(r)
	}
}

const errUnreadableInteger = "Could not read an integer!"

func readID() (int, err) {
	if _, e := skipWhitespace(); e != nil {
		return 0, newlineErr{errUnreadableInteger}
	}

	var idBuilder strings.Builder

	r, _, e := stdin.ReadRune()

	if e != nil || (r != '+' && r != '-' && !unicode.IsNumber(r)) {
		return 0, newlineErr{errUnreadableInteger}
	}

	for {
		idBuilder.WriteRune(r)
		r, _, e = stdin.ReadRune()

		if e != nil {
			return 0, newlineErr{errUnreadableInteger}
		} else if !unicode.IsNumber(r) {
			stdin.UnreadRune()

			break
		}
	}

	idStr := idBuilder.String()
	id, e := strconv.Atoi(idStr)

	if e != nil {
		return 0, newlineErr{errUnreadableInteger}
	}

	return id, nil
}

func skipWhitespace() (int, error) {
	numRead := 0

	for {
		r, _, e := stdin.ReadRune()

		if e != nil {
			return numRead, e
		} else if !unicode.IsSpace(r) {
			stdin.UnreadRune()

			return numRead, nil
		}

		numRead++
	}
}

func main() {
	os.Exit(start())
}
