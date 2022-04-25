package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
)

func finishItem(db *sql.DB) {
	n := getDeleteIndex(true)
	updateFinishItem(db, n)
}

func deleteItem(db *sql.DB) {
	n := getDeleteIndex(true)
	deleteItemDb(db, n)
}

func getDeleteIndex(finish bool) int {
	// Get correct usage string
	action := "delete"
	if finish {
		action = "finish"
	}

	// Check for arguments
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: wtodo %s <id>\n", action)
		os.Exit(0)
	}
	n, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid ID!\n")
		os.Exit(0)
	}

	return n
}
