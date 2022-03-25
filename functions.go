package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func list(todos []Item) {
	fmt.Printf("list\n")
}

func addItem(todos *[]Item) {

}

func editItem(todos *[]Item) {
	usageInfo := "Usage: wtodo " + os.Args[1] + " <id>"

	// Check for the ID command line argument
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, usageInfo)
		os.Exit(1)
	}

	// Find the element to edit
	key, _ := strconv.Atoi(os.Args[2])
	index := -1
	for i := range *todos {
		if (*todos)[i].Id == key {
			index = i
			break
		}
	}

	// Exit if not found
	if index == -1 {
		fmt.Fprintln(os.Stderr, "ID not found:", os.Args[2], "\n", usageInfo)
		list(*todos)
		os.Exit(1)
	}

	// TODO: Add editing for the other options

	// Create temp file to edit the name of the Item
	temp, err := os.CreateTemp("", "tmp")
	if err != nil {
		log.Fatal(err)
	}
	defer temp.Close()
	defer os.Remove(temp.Name())

	// Write previous name to the file
	temp.WriteString((*todos)[index].Name)

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
	(*todos)[index].Name = line
}

func finishItem(todos *[]Item) {

}

func deleteItem(todos *[]Item) {

}
