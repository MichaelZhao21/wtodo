package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Date mappings
var dateFormats = map[int]string{
	4:  "0102",
	5:  ":0304",
	8:  "01022006",
	9:  "0102-0304",
	13: "01022006-0304",
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
const LIGHT_GREEN_C = "\033[1;32m"
const LIGHT_YELLOW_C = "\033[1;33m"

func list(todos []Item) {
	// Filter list by done and not done
	notDone, _ := filterItems(todos)

	// Print header
	currDate := time.Now().Format("Monday January 2, 2006 (1/2/06) 3:04pm")
	fmt.Printf("%s===== %s%s%d Items To Do %s| %s%s%s%s =====%s\n", WHITE_C, RESET_C, CYAN_C, len(notDone), WHITE_C, RESET_C, YELLOW_C, currDate, WHITE_C, RESET_C)

	// If no items, print message and exit
	if len(notDone) == 0 {
		fmt.Printf("%sNothing left to do! Use %s%swtodo add%s%s to add more items.%s\n\n", WHITE_C, RESET_C, GREY_C, RESET_C, WHITE_C, RESET_C)
		return
	}

	// Filter each section by how far it is from due (<1 day, <1 week, other)
	today, soon, later := dateSortItems(notDone)

	// Print out all 3 sections
	if len(today) > 0 {
		fmt.Printf("%sDO TODAY%s\n\n", LIGHT_YELLOW_C, RESET_C)
		for _, t := range today {
			printItem(t, 1)
		}
	}

	if len(soon) > 0 {
		fmt.Printf("\n%sDO SOON%s\n\n", LIGHT_YELLOW_C, RESET_C)
		for _, t := range soon {
			printItem(t, 2)
		}
	}

	if len(later) > 0 {
		fmt.Printf("\n%sDO LATER (>1 week)%s\n\n", LIGHT_YELLOW_C, RESET_C)
		for _, t := range later {
			printItem(t, 3)
		}
	}
}

// Helper function to filter todos
func filterItems(todos []Item) (notDone []Item, done []Item) {
	for _, t := range todos {
		if t.Finished {
			done = append(done, t)
		} else {
			notDone = append(notDone, t)
		}
	}
	return notDone, done
}

// Filter items by how far they are from due and sort them
func dateSortItems(todos []Item) (today []Item, soon []Item, later []Item) {
	now := time.Now()

	// Sort based on due date
	sort.Slice(todos, func(p, q int) bool {
		return todos[p].Due.Before(todos[q].Due)
	})

	// Iterate through list and move to 3 lists
	todayEnd := time.Date(now.Year(), now.Month(), now.Day()+1, now.Hour(), now.Minute(), 1, 0, time.Local)
	soonEnd := time.Date(now.Year(), now.Month(), now.Day()+7, now.Hour(), now.Minute(), 1, 0, time.Local)
	for _, t := range todos {
		if t.Due.Before(todayEnd) {
			today = append(today, t)
		} else if t.Due.Before(soonEnd) {
			soon = append(soon, t)
		} else {
			later = append(later, t)
		}
	}

	// Return todos
	return today, soon, later
}

// Helper function to display one todo item
// Severity = 1 - red, 2 - yellow, 3 - green
func printItem(t Item, severity int) {
	due := t.Due.Format("Mon 1/2/06 3:04pm")
	dateCol := GREEN_C
	if severity == 1 {
		dateCol = RED_C
	} else if severity == 2 {
		dateCol = YELLOW_C
	}

	fmt.Printf("  %s%d. %s%s%s%s %s%s\n", DARK_GREY_C, t.Id, RESET_C, dateCol, due, WHITE_C, t.Name, RESET_C)
}

func editItem(todos *[]Item, nextId *int, add bool) {
	usageInfo := "Usage: wtodo " + os.Args[1] + " <id> [tags]"
	var temp Item
	var index int

	// Maintain different usage info and store into new object if adding an item
	// Otherwise, store the old object to edit in the temp variable
	if add {
		usageInfo = "Usage: wtodo " + os.Args[1] + "[tags]"
	} else {
		index = findItem(todos, usageInfo)
		temp = (*todos)[index]
	}

	// Get flags for edit command
	var p int
	var l, d, s, name string
	var n bool
	dateFormat := " | Formats: MMDDYYYY-HHmm, MMDD-HHmm, MMDDYYYY, MMDD, :HHmm ([M]onth, [D]ate, [Y]ear, [H]our, [m]inute) | Defaults: Today at 11:59pm"
	editFlags := flag.NewFlagSet("edit", flag.ExitOnError)
	editFlags.IntVar(&p, "p", -1, "Priority of the todo item | 1 - high, 2 - normal (default), 3 - low")
	editFlags.StringVar(&l, "l", "", "How long the task will take | [l]ong, [m]edium, [s]hort (default)")
	editFlags.StringVar(&d, "d", "", "Due date"+dateFormat)
	editFlags.StringVar(&s, "s", "", "Start date"+dateFormat)
	if add {
		editFlags.StringVar(&name, "n", "[New todo]", "Name of the todo item, Default: [New todo]")
	} else {
		editFlags.BoolVar(&n, "n", false, "Edit name, enable flag to use text editor to edit todo item name")
	}

	// Parse flags if there are any (no flags = new todo added with empty values)
	if !add {
		editFlags.Parse(os.Args[3:])
	} else if len(os.Args) > 2 {
		editFlags.Parse(os.Args[2:])
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
		switch l {
		case "s", "short":
			temp.Length = ShortTask
		case "m", "medium":
			temp.Length = MediumTask
		case "l", "long":
			temp.Length = LongTask
		default:
			fmt.Fprintln(os.Stderr, "Invalid Length:", l, "\nLength should be [l]ong, [m]edium, [s]hort")
			os.Exit(1)
		}
	}

	// Parse dates based on the avaliable formats
	if d != "" {
		temp.Due = parseDatetime(d, dateFormat)
	}

	if s != "" {
		temp.Start = parseDatetime(d, dateFormat)
	}

	// Edit name if tag enabled
	if n {
		temp.Name = editName(temp.Name)
	}

	// Add name
	if add {
		temp.Name = name
	}

	// If this is a new item, append it to the array and return
	// Otherwise replace the current item with the newly edited one
	if add {
		temp.Id = *nextId
		*todos = append(*todos, temp)
		*nextId++
	} else {
		(*todos)[index] = temp
	}
}

// Helper function to find an existing ID in the array
func findItem(todos *[]Item, usageInfo string) int {
	// If it is an edit, find the item id and replace it
	// Check for the ID command line argument
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, usageInfo)
		os.Exit(1)
	}

	// Find the element to edit
	key, _ := strconv.Atoi(os.Args[2])
	for i := range *todos {
		if (*todos)[i].Id == key {
			return i
		}
	}

	// Error if not found
	fmt.Fprintln(os.Stderr, "ID not found:", os.Args[2], "\n", usageInfo)
	list(*todos)
	os.Exit(1)
	return -1
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

// Helpfer function to parse dates
func parseDatetime(d string, dateFormat string) time.Time {
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
		fmt.Fprintln(os.Stderr, "Invalid due date:", d, dateFormat)
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

func finishItem(todos *[]Item) {

}

func deleteItem(todos *[]Item) {

}
