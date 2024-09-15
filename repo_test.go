package main

import (
	"testing"
)

func TestLoadIssues(t *testing.T) {
	repo := NewRepository("zachariahcox", "test")
	if repo == nil {
		t.Fail()
	}
}
