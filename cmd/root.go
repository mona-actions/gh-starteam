package cmd

import (
	"os"

	"github.com/mona-actions/gh-starteam/pkg/migrate"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "starteam",
	Short: "gh cli extension to assist in the migration of Starteam repositories to Github.",
	Long:  `gh cli extension to assist in the migration of Starteam repositories to Github.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// Get parameters
		historyFile, _ := cmd.Flags().GetString("history-file")
		repoPath, _ := cmd.Flags().GetString("repo-path")

		// Set ENV variables
		viper.Set("history-file", historyFile)
		viper.Set("repo-path", repoPath)

		// Execute the migration
		migrate.CreateGitRepo()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gh-starteam.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	// Initialize the configuration
	cobra.OnInitialize(initConfig)

	// Add flags for the history file
	rootCmd.Flags().StringP("history-file", "f", "", "Path to the history file")
	rootCmd.MarkFlagRequired("history-file")

	// Add flags for the repo path
	rootCmd.Flags().StringP("repo-path", "r", "", "Path to the repository")
	rootCmd.MarkFlagRequired("repo-path")
}

func initConfig() {
	// Configure Viper
	viper.SetEnvPrefix("GHST")
	viper.AutomaticEnv()
}
