package main

import (
	"regexp"
	"strings"
	"testing"
)

func assert(t *testing.T, should_be_true bool, message string) {
	if !should_be_true {
		t.Fatalf(message)
	}
}
func TestGH(t *testing.T) {
	args := []string{"version"}
	result := callCLI(args)
	want := regexp.MustCompile("$gh.*")
	assert(t, want.MatchString(result), strings.Join(args, " "))
}
