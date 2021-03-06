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
const YELLOW_C = "\033[33m" // More gold-like
const CYAN_C = "\033[36m"   // Really bright blue
const GREY_C = "\033[37m"
const WHITE_C = "\033[1;37m"
const DARK_GREY_C = "\033[1;30m"
const LIGHT_RED_C = "\033[1;31m"
const LIGHT_GREEN_C = "\033[1;32m"
const TITLE0_C = "\033[38;5;225m"
const TITLE1_C = "\033[38;5;159m"
const DATE0_C = "\033[38;5;124m"
const DATE1_C = "\033[38;5;203m"
const DATE2_C = "\033[38;5;222m"
const DATE3_C = "\033[38;5;192m"
const RATE0_C = "\033[38;5;157m"
const RATE1_C = "\033[38;5;229m"
const RATE2_C = "\033[38;5;215m"

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
	DbName   string
	Username string
}

// Command: wtodo <action> [tags] <text>
func main() {
	// Define list and main id incrementer
	var settings Settings
	var db *sql.DB

	// Load preferences from file
	loadPrefs(&settings)

	// If first time using system, generate a username
	setup(&settings, false)

	// If user wants to re-run the setup, do it before loading data
	if len(os.Args[1:]) > 0 && (os.Args[1] == "s" || os.Args[1] == "setup") {
		setup(&settings, true)
	}

	// Load data from database
	db = connectDb(settings)
	defer db.Close()

	// Case where there are no command line arguments
	if len(os.Args[1:]) < 1 {
		list(db)
		return
	}

	// Run commands based on the action statement
	switch os.Args[1] {
	case "list", "l", "setup", "s":
		list(db)
	case "add", "insert", "a", "i":
		editItem(db, true)
	case "edit", "e":
		editItem(db, false)
	case "finish", "f":
		finishItem(db)
	case "delete", "d":
		deleteItem(db)
	default:
		fmt.Printf("%sInvalid Action: %s\n%sUsage: wtodo <action> [options]\n", LIGHT_RED_C, os.Args[1], RESET_C)
		os.Exit(0)
	}
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
		}
	}

	// Save preferences if changed
	if noUser || dbSetup {
		savePrefs(settings)
	}
}
