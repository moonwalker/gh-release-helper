package main

import (
	"github.com/moonwalker/gh-release-helper/cmd"
)

var (
	version = "dev"
	commit  = "HEAD"
	date    = "n/a"
)

func main() {
	cmd.Execute(version, commit, date)
}
