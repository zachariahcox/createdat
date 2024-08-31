package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/cli/go-gh/v2"
	"golang.org/x/exp/constraints"
)

var batchSize = 25

func callCLI(cmd []string) string {
	stdout, stderr, err := gh.Exec(cmd...)
	if err != nil {
		log.Fatal(strings.Join(cmd, " "), "\n",
			stdout.String(), "\n",
			stderr.String(), "\n",
			err)
	}
	fmt.Println(stdout.String())
	return stdout.String()
}

func main() {
	// details to be pulled in from args
	owner := "slsa-framework"
	repo := "slsa"
	projectNumber := "5"

	// find Created Date field
	cmd := []string{"project", "field-list",
		"--owner", owner, projectNumber,
		"--format", "json",
		"--jq", ".fields[] | select(.name==\"Created Date\") | .id"}
	fieldId := callCLI(cmd)
	if fieldId == "" {
		log.Fatal("Could not find Created Date field")
	}
	fmt.Println("fieldId: " + fieldId)

	// fetch all issues in the project
	// load their creation dates and issue ids (these are different than the project "content" ids)
	cmd = []string{"issue", "list",
		"--repo", strings.Join([]string{owner, repo}, "/"),
		"--json", "createdAt,id,url",
		"--jq", ".[] | .id + \" \" + .createdAt + \" \" +.url"}
	id_date_pairs := strings.Split(callCLI(cmd), "\n")
	for _, p := range id_date_pairs {
		if p == "" {
			continue // sometimes these end with a newline
		}

		// super fancy parsing
		parts := strings.Split(p, " ")
		id := parts[0]
		date := parts[1]
		url := parts[2]
		fmt.Println("id: " + id + " date: " + date + " url: " + url)
	}

}

type ProjectItem struct {
	FieldIndex          int
	ProjectIndex        int
	ProjectId           string
	ProjectItemId       string
	FieldId             string
	ProjectV2FieldValue string // this is an https://docs.github.com/en/graphql/reference/input-objects#projectv2fieldvalue
}

func generateUpdateStatement(updates []ProjectItem) string {
	templateFile := "gql/update_issues.tmpl"
	t, err := template.New(templateFile).ParseFiles(templateFile)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, updates)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func makeChanges(updates []ProjectItem) {
	len_updates := len(updates)
	batchSize := (len_updates + batchSize - 1) / batchSize // this is just ceil, golang doesn't have int ceil???
	for i := 0; i < len_updates; i += batchSize {

		// golang HAS NO MIN FUNCTION FOR INTEGERS.
		end := i + batchSize
		if end > len_updates {
			end = len_updates
		}

		s := generateUpdateStatement(updates[i:end])
		args := []string{"api", "graphql", "--query", s}
		callCLI(args)
	}
}

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
