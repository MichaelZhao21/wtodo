package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
)

func finishItem(todos *[]Item, useDb bool, db *sql.DB) {
	// Check for arguments
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: wtodo finish <id>\n")
		os.Exit(0)
	}
	n, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid ID!\n")
		os.Exit(0)
	}

	// Either update in database or set to finished in list
	if useDb {
		updateFinishItem(db, n)
	} else {
		(*todos)[n].Finished = true
	}
}

func deleteItem(todos *[]Item) {

}
