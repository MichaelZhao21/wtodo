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
}

// Command: wtodo <action> [tags] <text>
func main() {
	// Define lists
	var todos []Item
	var finished []Item

	// Load data from file
	load(&todos, &finished)

	// Case where there are no command line arguments
	if len(os.Args[1:]) < 1 {
		list(todos, finished)
		return
	}

	// Run commands based on the action statement
	switch os.Args[1] {
	case "add", "insert", "a", "i":
		fmt.Printf("Add!\n")
	case "edit", "e":
		fmt.Printf("Edit!\n")
		edit()
	case "finish", "f":
		fmt.Printf("Finish!\n")
	case "delete", "d":
		fmt.Printf("Delete!\n")
	default:
		fmt.Fprintln(os.Stderr, "Invalid Action:", os.Args[1], "\nUsage: wtodo <action> [tags] <text>")
		os.Exit(1)
	}

	// Save data and exit
	save(&todos, &finished)
	os.Exit(0)
}
