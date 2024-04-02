package repo

import (
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/mona-actions/gh-starteam/internal/starteam"
	"github.com/spf13/viper"
)

type Commit struct {
	Author  string
	Date    time.Time
	Message string
	Files   []string
}

func CreateCommitsFromRevisions(revisions []starteam.Revision) []Commit {

	// Group revisions by date
	groupedRevisions := groupRevisionsByDate(revisions)

	// Create commits from grouped revisions
	var commits []Commit
	for _, revisions := range groupedRevisions {
		commit := Commit{
			Author:  revisions[0].Author,
			Date:    revisions[0].Date,
			Message: createMessageFromRevisions(revisions),
		}

		// Add files to commit
		for _, revision := range revisions {
			commit.Files = append(commit.Files, fixFilePath(revision.Folder, revision.FileName))
		}

		commits = append(commits, commit)
	}

	return commits
}

func groupRevisionsByDate(revisions []starteam.Revision) [][]starteam.Revision {
	groupedRevisions := make(map[string][]starteam.Revision)
	var dates []time.Time

	for _, revision := range revisions {
		date := revision.Date
		dateStr := date.String()

		if _, ok := groupedRevisions[dateStr]; !ok {
			// If the date is not in the map, add it to the slice
			dates = append(dates, date)
		}

		groupedRevisions[dateStr] = append(groupedRevisions[dateStr], revision)
	}

	// Sort the dates slice
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	// Create a slice of slices of revisions in date order
	var orderedRevisions [][]starteam.Revision
	for _, date := range dates {
		orderedRevisions = append(orderedRevisions, groupedRevisions[date.String()])
	}

	return orderedRevisions
}

func createMessageFromRevisions(revisions []starteam.Revision) string {
	message := ""

	for _, revision := range revisions {
		message += "Number: " + strconv.Itoa(revision.Number) + "\n"
		message += "Author: " + revision.Author + "\n"
		message += "Date: " + revision.Date.String() + "\n"
		message += "File: " + revision.FileName + "\n"
		message += "Folder: " + revision.Folder + "\n"
		message += "Memo: \n"
		message += revision.Memo + "\n"
		message += "-----------------------------\n"
	}

	return message
}

func fixFilePath(folder string, fileName string) string {
	// Remove "/home/dev/develop/C:/WorkingFolder/" from the folder
	folder = strings.Replace(folder, "/home/dev/develop/C:/WorkingFolder/", "", 1)

	// Remove "/home/dev/develop/" from the folder
	folder = strings.Replace(folder, "/home/dev/develop/", "", 1)

	// Remove "C:/WorkingFolder/" from the folder
	folder = strings.Replace(folder, "C:/WorkingFolder/", "", 1)

	// Add the repo-path to the beginning of the path and return
	return filepath.Join(viper.GetString("repo-path"), folder, fileName)
}
