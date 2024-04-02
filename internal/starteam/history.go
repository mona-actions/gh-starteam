package starteam

import (
	"log"
	"os"
	"sort"
	"strings"

	"github.com/spf13/viper"
)

func ProcessHistoryFile(histroyrFile string) []Revision {
	// Read the history file into memory
	data, err := os.ReadFile(viper.GetString("history-file"))
	if err != nil {
		log.Fatal(err)
	}

	// Convert the data into a slice of strings
	objects := strings.Split(string(data), "=============================================================================")

	// Remove the first three lines of the first object as they are not needed
	objects[0] = strings.Join(strings.Split(objects[0], "\n")[3:], "\n")

	// Remove the last object as it is just a list of folders
	objects = objects[:len(objects)-1]

	// Create a slice of files
	var revisions []Revision

	// Loop through the objects
	for _, object := range objects {
		// Process the object
		revisions = append(revisions, CreateRevisionsFromString(object)...)
	}

	// Loop through the files and fix missing folders
	var lastFolder string
	for i, revision := range revisions {
		if revision.Folder == "" {
			revisions[i].Folder = lastFolder
		} else {
			lastFolder = revision.Folder
		}
	}

	// Sort the revisions by date
	sort.Slice(revisions, func(i, j int) bool {
		return revisions[i].Date.Before(revisions[j].Date)
	})

	// return files
	return revisions
}
