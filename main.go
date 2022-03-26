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
		list(todos)
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
