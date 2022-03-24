package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
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
		log.Fatal(fmt.Sprint("Invalid Command: ", os.Args[1]))
	}

	// Save data and exit
	save(&todos, &finished)
	os.Exit(0)
}

func list(todos []Item, finished []Item) {
	fmt.Printf("list\n")
}

func edit() {
	// Create temp file
	temp, err := os.CreateTemp("", "tmp")
	if err != nil {
		log.Fatal(err)
	}
	defer temp.Close()
	defer os.Remove(temp.Name())

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
	content, err := ioutil.ReadFile(temp.Name())
	if err != nil {
		log.Fatal(err)
	}
	line := strings.Split(string(content), "\n")[0]
	fmt.Println(line)
}

// Load data from file
func load(todos *[]Item, finished *[]Item) {
	// Open data file
	f, err := os.OpenFile("wtodo.dat", os.O_RDONLY|os.O_CREATE, 0755)
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

	// Scan the scond line and split it
	// to get the length of 2 arrays
	scan.Scan()
	s = scan.Text()
	ss := strings.Split(s, " ")
	todosLen, _ := strconv.Atoi(ss[0])
	finishedLen, _ := strconv.Atoi(ss[1])

	// Iterate through both arrays and save all values
	for i := 0; i < todosLen; i++ {
		*todos = append(*todos, readItem(scan))
	}
	for i := 0; i < finishedLen; i++ {
		*finished = append(*finished, readItem(scan))
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
	rawDue, _ := strconv.ParseInt(ss[2], 10, 64)
	item.Due = time.Unix(rawDue, 0)
	rawStart, _ := strconv.ParseInt(ss[3], 10, 64)
	item.Start = time.Unix(rawStart, 0)

	// Read in line 2
	scan.Scan()
	s = scan.Text()
	item.Name = s

	// Return the item
	return item
}

// Save data to file
func save(todos *[]Item, finished *[]Item) {
	// Instantiate the stringbuilder
	sb := strings.Builder{}

	// Write the version on the first line
	sb.WriteString(Version)
	sb.WriteString("\n")

	// Save lengths of the arrays
	sb.WriteString(strconv.Itoa(len(*todos)))
	sb.WriteString(" ")
	sb.WriteString(strconv.Itoa(len(*finished)))
	sb.WriteString("\n")

	// Iterate through the todos and finished arrays and write all the data
	for _, item := range *todos {
		writeItem(item, &sb)
	}
	for _, item := range *finished {
		writeItem(item, &sb)
	}

	// Write all the data to the file
	os.WriteFile("wtodo.dat", []byte(sb.String()), fs.FileMode(os.O_TRUNC))
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
	(*sb).WriteString("\n")
	(*sb).WriteString(item.Name)
	(*sb).WriteString("\n")
}
