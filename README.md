# Micro Media Manager

`mediamanager` is a CLI application written in Go that allows users to create,
view, and manipulate a media collection. `mediamanager` is a rewrite of a
project from a course I took.

# Usage

[![asciicast](https://asciinema.org/a/p2dmOhvf5itZoFSiLukPW30Rv.svg)](https://asciinema.org/a/p2dmOhvf5itZoFSiLukPW30Rv)

## Introduction

* A "title" is a string with no leading or trailing whitespace and a single
  ASCII space character (`' '`) between non-whitespace characters.
* A "medium" is a string that contains no whitespace characters, such as "DVD"
  or "VHS".
* A "rating" is an integer between 1 and 5.
* An "ID" is a positive integer.
* A Record is a piece of media with a medium, title, rating, and unique ID
  assigned when it is created. Records are initially unrated.
* The Library is the set of Records.
* A Collection is a named subset of Records in the Library.
* The Catalog is the set of Collections.

## Input Processing

* Titles are read by consuming all input until the next newline, then stripping
  leading/trailing whitespace and compacting whitespace into a single space.
* Mediums are ready by skipping input until the next non-whitespace code point,
  then reading up until the next whitespace code point.
* Ratings are read as if by the regex `/\w*(?:+|-)?\d+/`. In other words,
  whitespace is skipped, an optional leading or trailing plus or minus sign is
  read, then a sequence of one or more numeric code points is read. Input stops
  at the first non-numeric code point.
* IDs are read the same way as ratings.
* Command strings are read character by character, skipping leading whitespace.

## Command Reference

* `fr <title>`: find Record. Find and print a Record in the Library, indexed by
  title.
* `pr <ID>`: print Record. Find and print a Record in the Library, indexed by
  ID.
* `pc <name>`: print Collection. Print a Collection in the Catalog.
* `pL`: print Library. Print all Records in the Library, sorted by title in
  ascending order.
* `pC`: print Catalog. Print all Collections in the Catalog, sorted by name in
  ascending order.
* `pa`: print allocations. Print the number of Records in the Library and the
  number of Collections in the Catalog.
* `ar <medium> <title>`: add Record. Add a new Record to the Library.
* `ac <name>`: add Collection. Add an empty Collection to the Catalog.
* `am <name> <ID>`: add member. Add a Record (indexed by ID) to a Collection.
* `mr <ID> <rating>`: modify rating. Change the rating of a Record.
* `dr <title>`: delete Record. Remove a Record from the Library.
* `dc <name>`: delete Collection. Remove a Collection from the Catalog.
* `dm <name> <ID>`: delete member. Remove a Record from a Collection.
* `cL`: clear Library. Remove all Records from the Library.
* `cC`: clear Catalog. Remove all Collections from the Catalog.
* `cA`: clear all. Clear the Library and the Catalog.
* `sA <filename>`: save all. Serialize the Library and Catalog to a file.
* `rA <filename>`: restore all. Deserialize the Library and Catalog from a file.
* `qq`: quit.
* `fs <string>`: find string. Print all Records that contain a substring,
  matching case insensitively.
* `lr`: list ratings. Print all Records in the Library, sorted by rating in
  descending order. Records with the same rating are sorted by title in
  ascending order.
* `cs`: Collection statistics. Print the number of Records that are a)
  contained in at least one Collection, b) contained in more than one
  Collection, and c) contained in Collections.
* `cc <firstSrcName> <secondSrcName> <dstName>`: combine Collections. Create a
  new Collection from the set union of two existing Collections, leaving the two
  source Collections unmodified.
* `mt <ID> <title>`: modify title. Change the title of a Record.

# License

`mediamanager` is licensed under the MIT license.
