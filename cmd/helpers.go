package cmd

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func gitrepo() string {
	p, err := run("git", "config", "--get", "remote.origin.url")
	if err != nil {
		log.Fatal(p)
	}

	res, err := run("basename", "-s", ".git", p)
	if err != nil {
		log.Fatal(res)
	}

	return res
}

func gitlog() string {
	t, err := run("git", "describe", "--tags", "--abbrev=0")
	if err != nil {
		t = "-10"
	} else {
		t = t + "..."
	}

	res, err := run("git", "log", "--pretty=format:- %H %s", t)
	if err != nil {
		log.Fatal(res)
	}

	return res
}

func run(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Env = os.Environ()

	out, err := cmd.CombinedOutput()
	res := strings.TrimSpace(string(out))

	return res, err
}
