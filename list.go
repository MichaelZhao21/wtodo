package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Function to list all items
func list(todos []Item, nextId int) {
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
	late, today, soon, later := dateSortItems(notDone)
	idWidth := strconv.Itoa(len(strconv.Itoa(nextId - 1)))

	// Print out all 4 sections
	if len(late) > 0 {
		fmt.Printf("%sOVERDUE%s\n", GREY_C, RESET_C)
		for _, t := range late {
			printItem(t, 0, idWidth)
		}
	}

	if len(today) > 0 {
		fmt.Printf("\n%sDO TODAY%s\n", GREY_C, RESET_C)
		for _, t := range today {
			printItem(t, 1, idWidth)
		}
	}

	if len(soon) > 0 {
		fmt.Printf("\n%sDO SOON%s\n", GREY_C, RESET_C)
		for _, t := range soon {
			printItem(t, 2, idWidth)
		}
	}

	if len(later) > 0 {
		fmt.Printf("\n%sDO LATER (>1 week)%s\n", GREY_C, RESET_C)
		for _, t := range later {
			printItem(t, 3, idWidth)
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
func dateSortItems(todos []Item) (late []Item, today []Item, soon []Item, later []Item) {
	now := time.Now()
	var never []Item

	// Sort based on due date
	sort.Slice(todos, func(p, q int) bool {
		return todos[p].Due.Before(todos[q].Due)
	})

	// Iterate through list and move to 3 lists
	todayEnd := time.Date(now.Year(), now.Month(), now.Day()+1, now.Hour(), now.Minute(), 1, 0, time.Local)
	soonEnd := time.Date(now.Year(), now.Month(), now.Day()+7, now.Hour(), now.Minute(), 1, 0, time.Local)
	for _, t := range todos {
		if t.Due.IsZero() {
			never = append(never, t)
		} else if t.Due.Before(now) {
			late = append(late, t)
		} else if t.Due.Before(todayEnd) {
			today = append(today, t)
		} else if t.Due.Before(soonEnd) {
			soon = append(soon, t)
		} else {
			later = append(later, t)
		}
	}

	// Sort never due items by ID and concat onto the end of later
	sort.Slice(todos, func(p, q int) bool {
		return todos[p].Id < todos[q].Id
	})
	later = append(later, never...)

	// Return todos
	return late, today, soon, later
}

// Helper function to display one todo item
// Severity = 0 - red bold, 1 - red, 2 - yellow, 3 - green
func printItem(t Item, severity int, idWidth string) {
	due := t.Due.Format("Mon 1/2/06 3:04pm")
	tags := strings.Join(t.Tags, ",")
	name := t.Name
	if len(name) > 30 {
		name = name[:27] + "..."
	}

	dueWidth := "21"
	nameWidth := "30"
	var dateCol string
	switch severity {
	case 0:
		dateCol = LIGHT_RED_C
	case 1:
		dateCol = RED_C
	case 2:
		dateCol = YELLOW_C
	default:
		if t.Due.IsZero() {
			dueWidth = "0"
			nameWidth = "50"
			due = ""
		}
		dateCol = GREEN_C
	}

	format := "  %s%" + idWidth + "d. %s%s%-" + dueWidth + "s%s%-" + nameWidth + "s%s %s%s%s\n"
	// println(format)
	fmt.Printf(format, DARK_GREY_C, t.Id, RESET_C, dateCol, due, WHITE_C, name, RESET_C, GREY_C, tags, RESET_C)
}
