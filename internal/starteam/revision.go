package starteam

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Revision struct {
	Author   string
	Date     time.Time
	Number   int
	FileName string
	Folder   string
	Memo     string
}

func CreateRevisionsFromString(data string) []Revision {
	// Split the metadata from the revision history
	re := regexp.MustCompile(`(?m)^----------------------------$`)
	splitData := re.Split(data, -1)

	// Separate out metadata from revision data
	metaData := splitData[0]
	revisionData := splitData[1:]

	// Parse the metadata
	metadata := parseMetadata(metaData)

	// Create a slice of revisions
	var revisions []Revision

	// Loop through the revision data
	for _, revision := range revisionData {
		// Parse the revision data
		revision := parseRevision(revision)

		// Add the metadata to the revision
		revision.Folder = metadata["folder"]
		revision.FileName = metadata["fileName"]

		// Append the revision to the slice
		revisions = append(revisions, revision)
	}

	// Return the slice of revisions
	return revisions
}

func parseMetadata(data string) map[string]string {
	// Split the metadata into lines
	lines := strings.Split(data, "\n")

	// Create a map to store the metadata
	metadata := make(map[string]string)

	// Loop through the lines
	for _, line := range lines {
		// If the line contains the string "Folder:", then extract the folder name and add it to any existing folder name
		if strings.Contains(line, "Folder:") {
			metadata["folder"] = strings.TrimSpace(strings.Split(line, "working dir: ")[1])
			// Remove last character
			metadata["folder"] = metadata["folder"][:len(metadata["folder"])-1]
		} else if strings.Contains(line, "History for:") {
			// If the line contains the string "History for:", then extract the file name
			metadata["fileName"] = strings.Replace(line, "History for: ", "", 1)
		}
	}

	// Return the metadata
	return metadata
}

func parseRevision(data string) Revision {
	// Create a new revision
	revision := Revision{}

	// Split the data into lines
	lines := strings.Split(data, "\n")

	// Loop through the lines
	for _, line := range lines {
		if strings.Contains(line, "Revision:") {
			// Extract the revision number
			numberStr := strings.Split(line, " ")[1]
			number, _ := strconv.Atoi(numberStr)
			revision.Number = number
		} else if strings.Contains(line, "Author:") && strings.Contains(line, "Date:") {
			// Find the start and end positions
			start := strings.Index(line, "Author:") + len("Author: ")
			end := strings.Index(line, "Date:")

			// Extract the author's name
			author := strings.TrimSpace(line[start:end])

			// Assign the author to the revision
			revision.Author = author

			// Extract the date
			dateStr := strings.TrimSpace(strings.Split(line, "Date: ")[1])

			// Parse the date
			date, err := time.Parse("1/2/06, 3:04:05 PM MST", dateStr)
			if err != nil {
				// Handle the error
				log.Fatal(err)
			}

			// Assign the date to the revision
			revision.Date = date
		} else if len(strings.TrimSpace(line)) == 0 {
			// Skip empty lines
		} else {
			// Add to the memo
			revision.Memo += line + "\n"
		}
	}

	// Return the revision
	return revision
}
