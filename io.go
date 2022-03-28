package main

import (
	"bufio"
	"io/fs"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
DATA FILE FORMAT:

======= | Header | =======
<version>
<length of todos> <nextId>
==========================

=============== | For each todo Item (3 lines) | ===============
<id> <length> <priority> <due date> <start date> <finished bool>
<name>
<comma seperated list of tags>
================================================================
*/

// Load data from file
func load(todos *[]Item, nextId *int) {
	path := getDataPath()

	// Open data file
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Create scanner
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanLines)

	// Read the first line in the file
	// If empty, then it is a new file!
	// Otherwise, it would be the version and we just ignore it
	scan.Scan()
	s := scan.Text()
	if len(s) == 0 {
		return
	}

	// Scan the scond line and split it to get the
	// length of the array and next id
	scan.Scan()
	s = scan.Text()
	ss := strings.Split(s, " ")
	todosLen, _ := strconv.Atoi(ss[0])
	*nextId, _ = strconv.Atoi(ss[1])

	// Iterate through both arrays and save all values
	for i := 0; i < todosLen; i++ {
		*todos = append(*todos, readItem(scan))
	}
}

func readItem(scan *bufio.Scanner) Item {
	// Create empty struct
	item := Item{}

	// Read in line 1
	scan.Scan()
	s := scan.Text()
	ss := strings.Split(s, " ")
	item.Id, _ = strconv.Atoi(ss[0])
	rawTask, _ := strconv.Atoi(ss[1])
	item.Length = TaskLength(rawTask)
	item.Priority, _ = strconv.Atoi(ss[2])
	rawDue, _ := strconv.ParseInt(ss[3], 10, 64)
	item.Due = time.Unix(rawDue, 0)
	rawStart, _ := strconv.ParseInt(ss[4], 10, 64)
	item.Start = time.Unix(rawStart, 0)
	item.Finished = ss[5] == "1"

	// Read in line 2
	scan.Scan()
	s = scan.Text()
	item.Name = s

	// Read in line 3
	scan.Scan()
	s = scan.Text()
	if s != "NULL" {
		item.Tags = strings.Split(s, ",")
	}

	// Return the item
	return item
}

// Save data to file
func save(todos *[]Item, nextId *int) {
	// Instantiate the stringbuilder
	sb := strings.Builder{}

	// Write the version on the first line
	sb.WriteString(Version)
	sb.WriteString("\n")

	// Save the length of the todos array and the next ID on the next line
	sb.WriteString(strconv.Itoa(len(*todos)))
	sb.WriteString(" ")
	sb.WriteString(strconv.Itoa(*nextId))
	sb.WriteString("\n")

	// Iterate through the todos array and write all the data
	for _, item := range *todos {
		writeItem(item, &sb)
	}

	// Write all the data to the file
	os.WriteFile(getDataPath(), []byte(sb.String()), fs.FileMode(os.O_TRUNC))
}

// Utility function to write a single Item
func writeItem(item Item, sb *strings.Builder) {
	(*sb).WriteString(strconv.Itoa(item.Id))
	(*sb).WriteString(" ")
	(*sb).WriteString(strconv.Itoa(int(item.Length)))
	(*sb).WriteString(" ")
	(*sb).WriteString(strconv.Itoa(item.Priority))
	(*sb).WriteString(" ")
	(*sb).WriteString(strconv.FormatInt(item.Due.Unix(), 10))
	(*sb).WriteString(" ")
	(*sb).WriteString(strconv.FormatInt(item.Start.Unix(), 10))
	(*sb).WriteString(" ")
	(*sb).WriteString(boolToString(item.Finished))
	(*sb).WriteString("\n")
	(*sb).WriteString(item.Name)
	(*sb).WriteString("\n")
	if len(item.Tags) > 0 {
		(*sb).WriteString(strings.Join(item.Tags, ","))
	} else {
		(*sb).WriteString("NULL")
	}
	(*sb).WriteString("\n")
}

func boolToString(in bool) string {
	if in {
		return "1"
	}
	return "0"
}

// Helper function to get the path of the data file
func getDataPath() string {
	// Get home file path and make data dir if not exists
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	path := dirname + "/.wtodo/wtodo.dat"
	os.Mkdir(dirname+"/.wtodo", fs.FileMode(0755))
	return path
}
