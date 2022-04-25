package main

import (
	"bufio"
	"io/fs"
	"log"
	"os"
	"strconv"
	"strings"
)

/*
PREFERENCES FILE FORMAT:

<version>
<username>
<use database (0/1 - rest of fields required)>
[db url (host)] [db port] [db username] [db password]
*/

// Loads the preferences from the user
func loadPrefs(settings *Settings) {
	path := getDataFilePath()

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
	// If empty, then it is a new file! (return and don't do anything)
	// Otherwise, it would be the version and we just ignore it
	scan.Scan()
	s := scan.Text()
	if len(s) == 0 {
		settings.UseDb = false
		return
	}

	// Read in the second line with the number data
	scan.Scan()
	settings.Username = scan.Text()

	// If this is false, don't continue
	scan.Scan()
	s = scan.Text()
	settings.UseDb = s == "1"
	if !settings.UseDb {
		return
	}

	// Read in the rest of the lines
	scan.Scan()
	s = scan.Text()
	ss := strings.Split(s, " ")
	settings.DbHost = ss[0]
	settings.DbPort, _ = strconv.Atoi(ss[1])
	settings.DbUser = ss[2]
	settings.DbPass = ss[3]
	settings.DbName = ss[4]
}

// Saves the preferences for the user
func savePrefs(settings *Settings) {
	// Instantiate the stringbuilder
	sb := strings.Builder{}

	// Write the version on the first line
	sb.WriteString(Version)
	sb.WriteString("\n")

	// Save the username and useDb boolean var on the next 2 lines
	sb.WriteString(settings.Username)
	sb.WriteString("\n")

	// If database is used, save the info for the database
	if settings.UseDb {
		sb.WriteString(settings.DbHost)
		sb.WriteString(" ")
		sb.WriteString(strconv.Itoa(settings.DbPort))
		sb.WriteString(" ")
		sb.WriteString(settings.DbUser)
		sb.WriteString(" ")
		sb.WriteString(settings.DbPass)
		sb.WriteString(" ")
		sb.WriteString(settings.DbName)
	}

	// Write all the data to the file
	os.WriteFile(getDataFilePath(), []byte(sb.String()), fs.FileMode(os.O_TRUNC))
}

// Helper function to get the path of the data file
func getDataFilePath() string {
	// Get home file path and make data dir if not exists
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	path := homeDir + "/.wtodo/prefs.dat"
	os.Mkdir(homeDir+"/.wtodo", fs.FileMode(0755))
	return path
}
