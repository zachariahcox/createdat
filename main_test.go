package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestParseUrl(t *testing.T) {
	for _, url := range []string{
		"https://github.com/orgs/repo_name/projects/1",
		"https://github.com/users/repo_name/projects/1",
		"github.com/orgs/repo_name/projects/1"} {
		scope, owner, projectNumber := parseUrl(url)
		if owner != "repo_name" {
			t.Fatal("owner should be repo_name but was", owner)
		}
		if projectNumber != "1" {
			t.Fatal("projectNumber should be '1' but was", projectNumber)
		}
		if scope != "org" && scope != "user" {
			t.Fatal("scope should be 'org' or 'user' but was", scope)
		}
	}
}

func TestNoProject(t *testing.T) {
	// the current exe is the built testing app -- recursively call it I guess?
	const RecursingToken = "TestNoProjectRecursing"
	if os.Getenv(RecursingToken) == RecursingToken {
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestNoProject")
	cmd.Env = append(os.Environ(), strings.Join([]string{RecursingToken, RecursingToken}, "="))
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatal("should fail with a usage print. instead got", err)
}
