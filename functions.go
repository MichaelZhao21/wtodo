package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

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
