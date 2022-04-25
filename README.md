# wtodo

A simple command line utility to keep track of your todos, with all data stored in a postgresql database.

## Installation

Simply clone the repo and install [`go`](https://go.dev/)!

## Execution and Build

To execute the code, change to the `wtodo` directory and run `go run .`

To build the script, change to the `wtodo` directory and run `go build -o wtodo`

## List of Commands

```
wtodo [l]ist - Lists out all todo items, also runs with no action specified, use -[c]ompleted to see all completed tasks
wtodo [a]dd - Create a new todo, type "wtodo add -h" for more options or no options for interactive prompt
wtodo [c]reate - Same as add
wtodo [e]dit - Edits a specific todo item
wtodo [f]inish - Marks an item as completed
wtodo [d]elete - Deletes a specific item
```
