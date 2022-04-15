package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"os/user"
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
const ORANGE_C = "\033[38;5;201m"

// 218

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

type Settings struct {
	UseDb    bool
	DbHost   string
	DbPort   int
	DbUser   string
	DbPass   string
	Username string
}

// Command: wtodo <action> [tags] <text>
func main() {
	// Define list and main id incrementer
	var todos []Item
	var nextId int
	var settings Settings
	var db *sql.DB

	// Load preferences from file
	loadPrefs(&settings)

	// If first time using system, generate a username
	setup(&settings, false)

	// Load data from file or database
	if settings.UseDb {
		db = connectDb(settings)
		defer db.Close()
	} else {
		loadFile(&todos, &nextId)
	}

	// Case where there are no command line arguments
	if len(os.Args[1:]) < 1 {
		list(todos, nextId, settings.UseDb, db)
		return
	}

	// Run commands based on the action statement
	switch os.Args[1] {
	case "setup", "s":
		setup(&settings, true)
	case "add", "insert", "a", "i":
		editItem(&todos, &nextId, settings, db, true)
	case "edit", "e":
		editItem(&todos, &nextId, settings, db, false)
	case "finish", "f":
		finishItem(&todos)
	case "delete", "d":
		deleteItem(&todos)
	default:
		fmt.Fprintln(os.Stderr, "Invalid Action:", os.Args[1], "\nUsage: wtodo <action> [options]")
		os.Exit(1)
	}

	// Save data and exit
	saveFile(&todos, &nextId)
	os.Exit(0)
}

// Generates a username if one is not already made
// Also ask for postgresql data
func setup(settings *Settings, dbSetup bool) {
	// If no username, generate one
	noUser := len(settings.Username) == 0
	if noUser {
		rand.Seed(time.Now().UnixNano())
		v := rand.Int()
		user, _ := user.Current()
		settings.Username = fmt.Sprintf("%s-%d", user.Username, v)
		fmt.Println(settings.Username)
		dbSetup = true
	}

	// If no database, setup one
	// Can also manually show this prompt
	if dbSetup {
		// Prompt user asking if they want to use a database
		// and get the info for the database if yes
		getDbInfo(settings)

		// If using database, connect and create tables
		if settings.UseDb {
			// Connect to the database
			db := connectDb(*settings)
			defer db.Close()

			// Create tables if not created
			createTables(db)

			// Show user data if setup
			list(nil, 0, settings.UseDb, db)
		}
	}

	if noUser || dbSetup {
		savePrefs(settings)
	}
}
