package main

import (
	"fmt"
	"os"
	"time"
)

// Just in case we have a version update that breaks the data file
// this version is written to the top of the data file to be compared
// Currently, this value is IGNORED
const Version = "1"

// Date mappings
var dateFormats = map[int]string{
	4:  "0102",
	5:  ":1504",
	8:  "01022006",
	9:  "0102-1504",
	13: "01022006-1504",
}

// Color mappings
const RESET_C = "\033[0m"
const RED_C = "\033[31m"
const GREEN_C = "\033[32m"
const YELLOW_C = "\033[33m"
const BLUE_C = "\033[34m"
const PURPLE_C = "\033[35m"
const CYAN_C = "\033[36m"
const GREY_C = "\033[37m"
const WHITE_C = "\033[1;37m"
const DARK_GREY_C = "\033[1;30m"
const LIGHT_RED_C = "\033[1;31m"
const LIGHT_GREEN_C = "\033[1;32m"
const LIGHT_YELLOW_C = "\033[1;33m"

type TaskLength int

const (
	ShortTask TaskLength = iota
	MediumTask
	LongTask
)

type Item struct {
	Id       int
	Name     string
	Due      time.Time
	Start    time.Time
	Length   TaskLength
	Priority int
	Finished bool
	Tags     []string
}

// Command: wtodo <action> [tags] <text>
func main() {
	// Define list and main id incrementer
	var todos []Item
	var nextId int

	// TEST CASES
	// todos = append(todos, Item{0, "This is the name of the todo", time.Now(), time.Time{}, ShortTask, 2, false})

	// Load data from file
	load(&todos, &nextId)

	// Case where there are no command line arguments
	if len(os.Args[1:]) < 1 {
		list(todos, nextId)
		return
	}

	// Run commands based on the action statement
	switch os.Args[1] {
	case "add", "insert", "a", "i":
		editItem(&todos, &nextId, true)
	case "edit", "e":
		editItem(&todos, &nextId, false)
	case "finish", "f":
		finishItem(&todos)
	case "delete", "d":
		deleteItem(&todos)
	default:
		fmt.Fprintln(os.Stderr, "Invalid Action:", os.Args[1], "\nUsage: wtodo <action> [options]")
		os.Exit(1)
	}

	// Save data and exit
	save(&todos, &nextId)
	os.Exit(0)
}
