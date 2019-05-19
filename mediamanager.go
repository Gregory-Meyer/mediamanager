package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

var stdin *bufio.Reader

func start() int {
	commands := map[string]func(*library, *catalog, *int) err{
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

func (c *catalog) String() string {
	collections := map[string](collection)(*c)

	if len(collections) == 0 {
		return "Catalog is empty"
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Catalog contains %d collections", len(collections)))

	for _, cat := range collections {
		builder.WriteRune('\n')
		builder.WriteString(fmt.Sprint(&cat))
	}

	return builder.String()
}

type collection struct {
	name    string
	members map[int]*record
}

func newCollection(name string) collection {
	return collection{name, make(map[int]*record)}
}

func (c *collection) String() string {
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

type record struct {
	medium         string
	title          string
	rating         int
	id             int
	numCollections int
}

func newRecord(medium, title string, id int) record {
	return record{medium, title, 0, id, 0}
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
	id, e := readInt()

	if e != nil {
		return e
	}

	if rec, ok := lib.byID[id]; ok {
		fmt.Println(rec)

		return nil
	}

	return newlineErr{"No record with that ID!"}
}

func printCollection(_ *library, cat *catalog, _ *int) err {
	name, _ := readWord()

	if col, ok := (*cat)[name]; ok {
		fmt.Println(&col)

		return nil
	}

	return newlineErr{"No collection with that name!"}
}

func printLibrary(lib *library, _ *catalog, _ *int) err {
	if len(lib.byTitle) == 0 {
		fmt.Println("Library is empty")

		return nil
	}

	titleSet := make([]string, 0, len(lib.byTitle))

	for title := range lib.byTitle {
		titleSet = append(titleSet, title)
	}

	sort.Strings(titleSet)

	fmt.Println("Library contains", len(titleSet), "records:")

	for _, title := range titleSet {
		r := lib.byTitle[title]
		fmt.Println(r)
	}

	return nil
}

func printCatalog(_ *library, cat *catalog, _ *int) err {
	fmt.Println(cat)

	return nil
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

	fmt.Println("Record", thisID, "added")

	return nil
}

func addCollection(_ *library, cat *catalog, id *int) err {
	name, _ := readWord()

	if _, ok := (*cat)[name]; ok {
		return newlineErr{"Catalog already has a collection with this name!"}
	}

	(*cat)[name] = newCollection(name)
	fmt.Sprintln("Collection", name, "added")

	return nil
}

func addMember(lib *library, cat *catalog, _ *int) err {
	name, _ := readWord()
	col, ok := (*cat)[name]

	if !ok {
		return newlineErr{"No collection with that name!"}
	}

	id, e := readInt()

	if e != nil {
		return e
	}

	// this complies with the spec, even though the error messages are out of
	// order, since they're orthogonal
	// there's no way a collection can have a member that doesn't exist in the
	// library, so this is safe
	if _, ok := col.members[id]; ok {
		return newlineErr{"Record is already a member in the collection!"}
	}

	rec, ok := lib.byID[id]

	if !ok {
		return newlineErr{"No record with that ID!"}
	}

	col.members[id] = rec
	rec.numCollections++
	fmt.Println("Member", id, rec.title, "added")

	return nil
}

func modifyRating(lib *library, _ *catalog, _ *int) err {
	id, e := readInt()

	if e != nil {
		return e
	}

	rec, ok := lib.byID[id]

	if !ok {
		return newlineErr{"No record with that ID!"}
	}

	newRating, e := readInt()

	if e != nil {
		return e
	}

	if newRating < 1 || newRating > 5 {
		return newlineErr{"Rating is out of range!"}
	}

	rec.rating = newRating
	fmt.Println("Rating for record", id, "changed to", newRating)

	return nil
}

func deleteRecord(lib *library, _ *catalog, _ *int) err {
	title, e := readTitle()

	if e != nil {
		return e
	}

	rec, ok := lib.byTitle[title]

	if !ok {
		return regularErr{"No record with that title!"}
	}

	if rec.numCollections > 0 {
		return regularErr{"Cannot delete a record that is a member of a collection!"}
	}

	id := rec.id

	delete(lib.byTitle, title)
	delete(lib.byID, id)

	fmt.Println("Record", id, title, "deleted")

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

func readInt() (int, err) {
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
