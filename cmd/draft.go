package cmd

import (
	"context"
	"log"
	"time"

	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

const (
	dateLayout = "2006-01-02"
)

var (
	relTarget, relName string
)

var draftCmd = &cobra.Command{
	Use:   "draft",
	Short: "Create or update a draft",

	Run: func(cmd *cobra.Command, args []string) {
		createUpdateDraft()
	},
}

func init() {
	draftCmd.Flags().StringVar(&relTarget, "target", "master", "Release target (commitish)")
	draftCmd.Flags().StringVar(&relName, "name", time.Now().Format(dateLayout), "Release name")
	rootCmd.AddCommand(draftCmd)
}

func createUpdateDraft() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	var err error
	var rel *github.RepositoryRelease
	rels, _, _ := client.Repositories.ListReleases(ctx, ghOwner, ghRepo, nil)
	if rels != nil {
		for _, r := range rels {
			if *r.Draft {
				rel = r
				break
			}
		}
	}

	releaseBody := gitlog()
	if rel != nil {
		rel.Name = &relName
		rel.Body = &releaseBody
		rel.CreatedAt = &github.Timestamp{time.Now()}
		_, _, err = client.Repositories.EditRelease(ctx, ghOwner, ghRepo, *rel.ID, rel)
		if err != nil {
			log.Fatal("Failed to edit release", err)
		}
	} else {
		tagName := ""
		draft := true
		rel = &github.RepositoryRelease{
			TagName:         &tagName,
			TargetCommitish: &relTarget,
			Name:            &relName,
			Body:            &releaseBody,
			Draft:           &draft,
		}
		_, _, err = client.Repositories.CreateRelease(ctx, ghOwner, ghRepo, rel)
		if err != nil {
			log.Fatal("Failed to create release", err)
		}
	}
}
