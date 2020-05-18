package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-github/v31/github"
	"golang.org/x/oauth2"
)

var (
	token  = flag.String("token", "$GITHUB_TOKEN", "github access token")
	owner  = flag.String("owner", "moonwalker", "github owner (org or user)")
	repo   = flag.String("repo", "", "github repository")
	target = flag.String("target", "master", "release target commitish")
	name   = flag.String("name", "", "release target commitish")
)

const (
	dateLayout = "2006-01-02"
)

func init() {
	flag.Parse()
}

func main() {
	if token == nil || len(*token) == 0 {
		*token = os.Getenv("GITHUB_TOKEN")
	}
	if token == nil || len(*token) == 0 {
		log.Fatal(fmt.Errorf("github token not specified"))
	}

	if repo == nil || len(*repo) == 0 {
		log.Fatal(fmt.Errorf("github repo not specified"))
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	rels, _, err := client.Repositories.ListReleases(ctx, *owner, *repo, nil)
	if err != nil {
		log.Fatal(err)
	}

	var rel *github.RepositoryRelease
	for _, r := range rels {
		if *r.Draft {
			rel = r
			break
		}
	}

	releaseBody := gitlog()
	if name == nil || len(*name) == 0 {
		*name = time.Now().Format(dateLayout)
	}

	if rel != nil {
		rel.Name = name
		rel.Body = &releaseBody
		rel.CreatedAt = &github.Timestamp{time.Now()}
		_, _, err = client.Repositories.EditRelease(ctx, *owner, *repo, *rel.ID, rel)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		tagName := ""
		draft := true
		rel = &github.RepositoryRelease{
			TagName:         &tagName,
			TargetCommitish: target,
			Name:            name,
			Body:            &releaseBody,
			Draft:           &draft,
		}
		_, _, err = client.Repositories.CreateRelease(ctx, *owner, *repo, rel)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func gitlog() string {
	tag := run("git", "describe", "--tags", "--abbrev=0")
	return run("git", "log", "--pretty=format:- %H %s", tag+"...")
}

func run(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	cmd.Env = os.Environ()

	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Dir = filepath.Dir(ex)

	out, err := cmd.CombinedOutput()
	res := strings.TrimSpace(string(out))
	if err != nil {
		log.Fatal(res)
	}

	return res
}
