package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	name = "ghr"
)

var (
	ghToken, ghOwner, ghRepo string
)

var rootCmd = &cobra.Command{
	Use:   name,
	Short: "GitHub release helper",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&ghToken, "token", "$GITHUB_TOKEN", "GitHub access token")
	rootCmd.PersistentFlags().StringVar(&ghOwner, "owner", "moonwalker", "GitHub owner org or user")
	rootCmd.PersistentFlags().StringVar(&ghRepo, "repo", "", "GitHub repository")
}

func Execute(version, commit, date string) {
	verInfo = &versionInfo{version, commit, date}

	if ghToken == "$GITHUB_TOKEN" {
		ghToken = os.Getenv("GITHUB_TOKEN")
	}

	if len(ghRepo) == 0 {
		ghRepo = gitrepo()
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
