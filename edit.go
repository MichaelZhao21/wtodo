package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Function to edit and add items
func editItem(db *sql.DB, add bool) {
	usageInfo := "Usage: wtodo " + os.Args[1] + " <id> [tags]"
	var oldItem Item

	// Create and set default temp values
	temp := Item{}
	temp.Length = ShortTask
	temp.Priority = 2

	// Maintain different usage info and store into new object if adding an item
	// Otherwise, store the old object to edit in the temp variable if not using a database
	if add {
		usageInfo = "Usage: wtodo " + os.Args[1] + "[tags]"
	} else {
		oldItem = findItem(usageInfo, db)
	}

	// Get flags for edit command
	var p int
	var l, d, s, name, t string
	var n bool
	dateFormatSimple := "MMDDYYYY-HHmm, MMDD-HHmm, MMDDYYYY, MMDD, :HHmm, 0"
	dateFormat := "Formats: MMDDYYYY-HHmm, MMDD-HHmm, MMDDYYYY, MMDD, :HHmm, 0 ([M]onth, [D]ate, [Y]ear, [H]our, [m]inute, 0=none) | Defaults: Today at 11:59pm"
	editFlags := flag.NewFlagSet("add/edit", flag.ExitOnError)
	editFlags.IntVar(&p, "p", -1, "Priority of the todo item | 3 - high, 2 - normal (default), 1 - low")
	editFlags.StringVar(&l, "l", "", "How long the task will take | [l]ong, [m]edium, [s]hort (default)")
	editFlags.StringVar(&d, "d", "", "Due date | "+dateFormat)
	editFlags.StringVar(&s, "s", "", "Start date | "+dateFormat)
	editFlags.BoolVar(&n, "en", false, "Edit name (only used if editing), enable flag to use text editor to edit todo item name")
	editFlags.StringVar(&name, "n", "", "Name of the todo item, REQUIRED")
	editFlags.StringVar(&t, "t", "", "Tags (Comma-seperated)")

	// Parse flags if there are any
	if !add {
		editFlags.Parse(os.Args[3:])
	} else if len(os.Args) > 2 {
		editFlags.Parse(os.Args[2:])
	} else {
		// If there are no flags, run the interactive builder
		interactiveAdd(&temp, dateFormatSimple)
	}

	// Edit priority and check for the correct range of numbers
	if p != -1 {
		if p > 0 && p <= 3 {
			temp.Priority = p
		} else {
			fmt.Fprintln(os.Stderr, "Invalid Priority:", p, "\nPriority should be (1-3): 1 - high, 2 - normal, 3 - low")
			os.Exit(1)
		}
	}

	// Edit length and check for the correct values
	if l != "" {
		temp.Length = parseLength(l)
	}

	// Parse dates based on the avaliable formats
	if d != "" {
		temp.Due = parseDatetime(d, dateFormat)
	}

	if s != "" {
		temp.Start = parseDatetime(d, dateFormat)
	}

	// Edit name if tag enabled
	if !add && n {
		temp.Name = editName(oldItem.Name)
	}

	// Edit tags if valid and not empty
	if t != "" {
		temp.Tags = strings.Split(t, ",")
	}

	// Name field is required for adding a todo
	if len(os.Args) > 2 {
		if add && name == "" {
			fmt.Fprintln(os.Stderr, "Name field (-n) is required!")
			os.Exit(1)
		} else if name != "" {
			temp.Name = name
		}
	}

	// Add or update from database
	if add {
		insertItem(db, temp)
	} else {
		updateItem(db, temp)
	}
}

// Helper function to find an existing item in the database
func findItem(usageInfo string, db *sql.DB) Item {
	// If it is an edit, find the item id and replace it
	// Check for the ID command line argument
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, usageInfo)
		os.Exit(1)
	}

	// Select from database
	key, _ := strconv.Atoi(os.Args[2])
	item := selectItem(db, key)
	// TODO: Add error if item is not found
	return item

	// // Error if not found
	// log.Fatalln("ID not found:", os.Args[2], "\n", usageInfo)
	// return 0, Item{}
}

// Helper function to edit the name of an Item using the default text editor
func editName(oldName string) string {
	// Create temp file to edit the name of the Item
	temp, err := os.CreateTemp("", "tmp")
	if err != nil {
		log.Fatal(err)
	}
	defer temp.Close()
	defer os.Remove(temp.Name())

	// Write previous name to the file
	temp.WriteString(oldName)

	// Use $EDITOR as default editor or use vim
	editor := os.Getenv("EDITOR")
	if len(editor) == 0 {
		fmt.Printf("No editor set in env $EDITOR, using vim as default\n")
		editor = "vim"
	}

	// Run the command using standard IO
	cmd := exec.Command(editor, temp.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	// Read the file that the user wrote to
	// and save that data to the array
	content, err := ioutil.ReadFile(temp.Name())
	if err != nil {
		log.Fatal(err)
	}
	line := strings.Split(string(content), "\n")[0]
	return line
}

// Helper function to parse the string length to a TaskLength enum
func parseLength(l string) TaskLength {
	switch l {
	case "s", "short":
		return ShortTask
	case "m", "medium":
		return MediumTask
	case "l", "long":
		return LongTask
	default:
		fmt.Fprintln(os.Stderr, "Invalid Length:", l, "\nLength should be [l]ong, [m]edium, [s]hort")
		os.Exit(1)
	}
	return ShortTask
}

// Helper function to parse dates
func parseDatetime(d string, dateFormat string) time.Time {
	// Return zero time if string empty
	if d == "" || d == "0" {
		return time.Time{}
	}

	// Set defaults
	now := time.Now()
	year := now.Year()
	month := now.Month()
	day := now.Day()
	hour := 23
	minute := 59

	// Parse time form input string
	parsed, err := time.Parse(dateFormats[len(d)], d)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid due date:", d, "|", dateFormat)
		os.Exit(1)
	}

	// Modify defaults based on input string
	switch len(d) {
	case 8: // MMDDYYYY
		year = parsed.Year()
		month = parsed.Month()
		day = parsed.Day()
	case 5: // :HHmm
		hour = parsed.Hour()
		minute = parsed.Minute()
	case 13: // MMDDYYYY-HHmm
		year = parsed.Year()
		fallthrough
	case 9: // MMDD-HHmm
		hour = parsed.Hour()
		minute = parsed.Minute()
		fallthrough
	case 4:
		month = parsed.Month()
		day = parsed.Day()
	}

	// Return newly created date
	return time.Date(year, month, day, hour, minute, 0, 0, time.Local)
}

func interactiveAdd(todo *Item, dateFormat string) {
	read := bufio.NewReader(os.Stdin)
	for todo.Name == "" {
		fmt.Printf("%sEnter name %s[Required]%s ", YELLOW_C, GREY_C, RESET_C)
		name, _ := read.ReadString('\n')
		todo.Name = strings.Trim(name, " \n")
	}

	fmt.Printf("%sEnter priority %s(1 high, 2 normal, 3 low) [Default: 2]%s ", YELLOW_C, GREY_C, RESET_C)
	priority, _ := read.ReadString('\n')
	todo.Priority, _ = strconv.Atoi(priority[:len(priority)-1])
	if priority == "\n" {
		todo.Priority = 2
	} else if todo.Priority < 1 || todo.Priority > 3 {
		fmt.Printf("%sInvalid priority %d, defaulting to normal%s\n", RED_C, todo.Priority, RESET_C)
		todo.Priority = 2
	}

	fmt.Printf("%sEnter task length %s([s]hort, [m]edium, [l]ong) [Default: short]%s ", YELLOW_C, GREY_C, RESET_C)
	l, _ := read.ReadString('\n')
	if l == "\n" {
		todo.Length = ShortTask
	} else {
		todo.Length = parseLength(strings.Trim(l, "\n "))
	}

	fmt.Printf("%sEnter due date %s(%s) [Default: none]%s ", YELLOW_C, GREY_C, dateFormat, RESET_C)
	d, _ := read.ReadString('\n')
	todo.Due = parseDatetime(d[:len(d)-1], dateFormat)

	fmt.Printf("%sEnter start date %s(%s) [Default: none]%s ", YELLOW_C, GREY_C, dateFormat, RESET_C)
	s, _ := read.ReadString('\n')
	todo.Start = parseDatetime(s[:len(s)-1], dateFormat)
}
