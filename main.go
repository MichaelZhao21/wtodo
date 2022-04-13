package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/user"
	"strconv"
	"strings"
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

	// Load preferences from file
	loadPrefs(&settings)

	// If first time using system, generate a username
	setup(&settings, false)

	// Load data from file or database
	if settings.UseDb {
		db := connectDb(settings)
		defer db.Close()
		loadDb(&todos, settings)
	} else {
		loadFile(&todos, &nextId)
	}

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
	if dbSetup {
		// Prompt user asking if they want to use a database
		// and get the info for the database if yes
		getDbInfo(settings)

		// Connect to the database
		db := connectDb(*settings)
		defer db.Close()

		// Create tables if not created
		db.Exec("CREATE TABLE IF NOT EXISTS Item (id integer PRIMARY KEY, name varchar(100) NOT NULL, due timestamp with time zone, start timestamp with time zone, length smallint, priority smallint, finished boolean)")
		db.Exec("CREATE TABLE IF NOT EXISTS Tag (item_id integer PRIMARY KEY, name varchar(50))")
	}

	if noUser || dbSetup {
		savePrefs(settings)
	}
}

func getDbInfo(settings *Settings) {
	// Reset settings
	settings.DbHost = ""
	settings.DbPort = 0
	settings.DbUser = ""
	settings.DbPass = ""
	settings.UseDb = false

	// Create buffer reader
	read := bufio.NewReader(os.Stdin)

	// Ask the user if they want to use the database
	fmt.Printf("%sUse Postgresql database? (y/n) [Default n - use local data file]:%s ", YELLOW_C, RESET_C)
	useDb, _ := read.ReadString('\n')
	settings.UseDb = strings.ToLower(strings.Trim(useDb, "\n")) == "y"

	// If no database we done
	if !settings.UseDb {
		return
	}

	// Prompt for the rest of the inforamtion for the database
	fmt.Printf("%s\nDatabase Setup\n==============\n%s%sMake sure you have created a database named \"wtodo\" as we will connect to that database!\n%s", WHITE_C, RESET_C, LIGHT_GREEN_C, RESET_C)
	for settings.DbHost == "" {
		fmt.Printf("%s> Enter Database Host:%s ", YELLOW_C, RESET_C)
		name, _ := read.ReadString('\n')
		settings.DbHost = strings.Trim(name, " \n")
	}
	for settings.DbPort == 0 {
		fmt.Printf("%s> Enter Database Port:%s ", YELLOW_C, RESET_C)
		name, _ := read.ReadString('\n')
		settings.DbPort, _ = strconv.Atoi(strings.Trim(name, " \n"))
	}
	for settings.DbUser == "" {
		fmt.Printf("%s> Enter Database Username:%s ", YELLOW_C, RESET_C)
		name, _ := read.ReadString('\n')
		settings.DbUser = strings.Trim(name, " \n")
	}
	for settings.DbPass == "" {
		fmt.Printf("%s> Enter Database Password:%s ", YELLOW_C, RESET_C)
		name, _ := read.ReadString('\n')
		settings.DbPass = strings.Trim(name, " \n")
	}
	fmt.Printf("%s\nSetup complete!%s\n===============\nHost: %s\nPort: %d\nUsername: %s\nPassword: %s\n\n", WHITE_C, RESET_C, settings.DbHost, settings.DbPort, settings.DbUser, settings.DbPass)
}
