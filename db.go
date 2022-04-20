package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func getDbInfo(settings *Settings) {
	// Reset settings
	settings.DbHost = ""
	settings.DbPort = 0
	settings.DbUser = ""
	settings.DbPass = ""
	settings.UseDb = false
	settings.DbName = ""

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
	fmt.Printf("%s\nDatabase Setup\n==============\n%s", LIGHT_GREEN_C, RESET_C)
	for settings.DbHost == "" {
		fmt.Printf("%s> Enter Database Host:%s ", YELLOW_C, RESET_C)
		name, _ := read.ReadString('\n')
		settings.DbHost = strings.Trim(name, " \n")
	}

	fmt.Printf("%s> Enter Database Port [Default 5432]:%s ", YELLOW_C, RESET_C)
	port, _ := read.ReadString('\n')
	var err error
	settings.DbPort, err = strconv.Atoi(strings.Trim(port, " \n"))
	if err != nil {
		settings.DbPort = 5432
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

	for settings.DbName == "" {
		fmt.Printf("%s> Enter Database Name:%s ", YELLOW_C, RESET_C)
		name, _ := read.ReadString('\n')
		settings.DbName = strings.Trim(name, " \n")
	}

	fmt.Printf("%s\nSetup complete!%s\n===============\nHost: %s\nPort: %d\nUsername: %s\nPassword: %s\nDB Name: %s\n\n", WHITE_C, RESET_C, settings.DbHost, settings.DbPort, settings.DbUser, settings.DbPass, settings.DbName)
}

// Connects to database using the info stored in settings
func connectDb(settings Settings) *sql.DB {
	// Connect to the database
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", settings.DbHost, settings.DbPort, settings.DbUser, settings.DbPass, settings.DbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Could not open DB: ", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Database ping failed: ", err)
	}
	return db
}

// Create database tables
func createTables(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal("Database disconnected :( ", err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS Item (id serial PRIMARY KEY, name varchar(100) NOT NULL, due timestamp with time zone, start timestamp with time zone, length smallint, priority smallint, finished boolean);")
	_, err2 := db.Exec("CREATE TABLE IF NOT EXISTS Tag (item_id integer PRIMARY KEY, name varchar(50));")
	if err != nil {
		log.Fatal("Error creating item table:", err)
	}
	if err2 != nil {
		log.Fatal("Error creating tag table:", err2)
	}
}

// Selects all from item table
func selectAll(db *sql.DB) []Item {
	// Load current timezone
	americaTime := time.Now().Location()

	// Perform select
	q := `SELECT * FROM Item`
	rows, err := db.Query(q)
	if err != nil {
		log.Fatal("Error selecting all items:", err)
	}

	// Iterate through selection and save to struct
	var temp []Item
	for rows.Next() {
		var it Item
		err := rows.Scan(&it.Id, &it.Name, &it.Due, &it.Start, &it.Length, &it.Priority, &it.Finished)
		if err != nil {
			log.Fatal("Error loading row:", err)
		}
		it.Due = it.Due.In(americaTime)
		temp = append(temp, it)
	}

	return temp
}

// Insert item into database
func insertItem(db *sql.DB, item Item) {
	_, err := db.Exec("INSERT INTO Item VALUES (DEFAULT, $1, $2, $3, $4, $5, $6)", item.Name, item.Due, item.Start, item.Length, item.Priority, item.Finished)
	if err != nil {
		panic(err.Error())
	}
}

// Update item from database
func updateItem(db *sql.DB, item Item) {
	_, err := db.Exec("UPDATE Item SET name=$1, due=$2, start=$3, length=$4, priority=$5, finished=$6 WHERE id=$7", item.Name, item.Due, item.Start, item.Length, item.Priority, item.Finished, item.Id)
	if err != nil {
		panic(err.Error())
	}
}

// Select specific item from database
func selectItem(db *sql.DB, key int) Item {
	rows, err := db.Query("SELECT * FROM Item WHERE id=$1", key)
	if err != nil {
		panic(err.Error())
	}

	var temp Item
	if rows.Next() {
		err = rows.Scan(&temp.Id, &temp.Name, &temp.Due, &temp.Start, &temp.Length, &temp.Priority, &temp.Finished)
		if err != nil {
			panic(err.Error())
		}
	}

	return temp
}
