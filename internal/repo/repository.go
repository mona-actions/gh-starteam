package repo

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/mona-actions/gh-starteam/internal/starteam"
	"github.com/spf13/viper"
)

func CreateGitRepo(revisions []starteam.Revision) {
	// Get the parameters
	repoPath := viper.GetString("repo-path")

	// Init the new repository
	// fs := osfs.New(repoPath)
	// storer := filesystem.NewStorage(fs, nil)
	gitRepo, _ := git.PlainInit(repoPath, false)
	workingTree, _ := gitRepo.Worktree()

	// Create commits from revisions
	commits := CreateCommitsFromRevisions(revisions)

	// Commit counter
	commitCounter := 0

	// Loop through the commits and create them in the repository
	for _, commit := range commits {
		// Loop through the files and create them in the repository
		for _, file := range commit.Files {
			createFile(file, commit.Message)
		}

		// Add all the files to the working tree
		_, err := workingTree.Add(".")
		if err != nil {
			panic(err)
		}

		// Commit the files
		_, err = workingTree.Commit(commit.Message, &git.CommitOptions{
			Author: &object.Signature{
				Name:  commit.Author,
				Email: "test@test.com",
				When:  commit.Date,
			},
		})
		if err != nil {
			panic(err)
		}

		// Print the progress
		commitCounter++
		fmt.Println("Commits Completed: " + strconv.Itoa(commitCounter) + "/" + strconv.Itoa(len(commits)))
	}
}

func createFile(filePath string, content string) {
	// Check if the file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// Create the directories leading up to the file
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			panic(err)
		}

		// Create the file
		file, err := os.Create(filePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// Write the content to the file
		_, err = file.WriteString(content)
		if err != nil {
			panic(err)
		}
	}
}
