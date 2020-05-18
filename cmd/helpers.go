package cmd

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func gitrepo() string {
	p := run("git", "rev-parse", "--show-toplevel")
	return run("basename", p)
}

func gitlog() string {
	tag := run("git", "describe", "--tags", "--abbrev=0")
	return run("git", "log", "--pretty=format:- %H %s", tag+"...")
}

func run(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	cmd.Env = os.Environ()

	out, err := cmd.CombinedOutput()
	res := strings.TrimSpace(string(out))
	if err != nil {
		log.Fatal(res)
	}

	return res
}
