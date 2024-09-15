package main

import (
	"testing"
)

func TestLoadIssues(t *testing.T) {
	repo := &Repository{
		Owner: "zachariahcox",
		Name:  "test",
	}

	repo.LoadIssues()
}
