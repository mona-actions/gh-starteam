package migrate

import (
	"github.com/mona-actions/gh-starteam/internal/repo"
	"github.com/mona-actions/gh-starteam/internal/starteam"

	"github.com/spf13/viper"
)

func CreateGitRepo() {
	// Get the parameters
	historyFile := viper.GetString("history-file")

	// Process History File
	revisions := starteam.ProcessHistoryFile(historyFile)

	// Creating Git Repo
	repo.CreateGitRepo(revisions)
}
