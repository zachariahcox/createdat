package main

import (
	"testing"
)

func TestParseUrl(t *testing.T) {
	for _, url := range []string{
		"https://github.com/orgs/repo_name/projects/1",
		"github.com/orgs/repo_name/projects/1"} {
		scope, owner, projectNumber := parseUrl(url)
		if owner != "repo_name" {
			t.Fatal("owner should be repo_name but was", owner)
		}
		if projectNumber != "1" {
			t.Fatal("projectNumber should be '1' but was", projectNumber)
		}
		if scope != "orgs" {
			t.Fatal("scope should be 'orgs' but was", scope)
		}
	}
}
