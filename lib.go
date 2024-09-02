package main

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/cli/go-gh/v2"
)

var batchSize = 25
var DEBUG = true

//go:embed gql/*
var GqlFiles embed.FS

func callCLI(cmd []string) string {
	stdout, stderr, err := gh.Exec(cmd...)
	if err != nil {
		log.Fatal(strings.Join(cmd, " "), "\n",
			stdout.String(), "\n",
			stderr.String(), "\n",
			err)
	}
	return stdout.String()
}

func loadTemplate(filePath string) *template.Template {
	templateFile := filePath
	t, err := template.New(templateFile).ParseFiles(templateFile)
	if err != nil {
		panic(err)
	}
	return t
}

type Project struct {
	id       string
	title    string
	owner    string
	repo     string
	number   string
	items    *[]*ProjectItem
	contents *[]*ProjectItem
}

func NewProject(owner string, repo string, number string) *Project {
	p := new(Project)
	p.owner = owner
	p.repo = repo
	p.number = number
	p.contents = getProjectContents(p)
	return p
}

type Content struct {
	id        string
	url       string
	number    float64
	title     string
	closed    bool
	createdAt string
	closedAt  string
}
type PI struct {
	id          string
	createdAt   string
	content     Content
	contentType string
}

type ProjectItem struct {
	FieldIndex          int
	ProjectIndex        int
	ProjectId           string
	ProjectItemId       string
	FieldId             string
	ProjectV2FieldValue string // this is an https://docs.github.com/en/graphql/reference/input-objects#projectv2fieldvalue
}
type Issue struct {
	url  string
	id   string
	date string
}

func (issue *Issue) String() string {
	return "id: " + issue.id + " date: " + issue.date + " url: " + issue.url
}

func getFieldId(p *Project, fieldName string) string {
	cmd := []string{"project", "field-list",
		"--owner", p.owner, p.number,
		"--format", "json",
		"--jq", ".fields[] | select(.name==\"" + fieldName + "\") | .id"}
	fieldId := callCLI(cmd)
	if fieldId == "" {
		log.Fatal("Could not find field named '" + fieldName + "'")
	}
	return fieldId
}

func getIssues(p *Project) *[]*Issue {

	fields := []string{"createdAt", "id", "url"}
	d := make([]string, 0, len(fields))
	for _, s := range fields {
		d = append(d, "."+s)
	}
	cmd := []string{"issue", "list",
		"--repo", strings.Join([]string{p.owner, p.repo}, "/"),
		"--json", strings.Join(fields, ","),
		"--jq", ".[] | " + strings.Join(d, "+\" \"+")}

	// super fancy parsing
	fmt.Println(strings.Join(cmd, " "))
	id_date_pairs := strings.Split(callCLI(cmd), "\n")
	issues := make([]*Issue, 0, len(id_date_pairs))

	for _, p := range id_date_pairs {
		if p == "" {
			continue // sometimes these end with a newline
		}
		parts := strings.Split(p, " ")
		issue := Issue{date: parts[0], id: parts[1], url: parts[2]}
		issues = append(issues, &issue)
	}
	return &issues
}

func getProjectContents(p *Project) *[]*ProjectItem {
	bytes, err := GqlFiles.ReadFile("gql/get_project_contents.gql")
	if err != nil {
		panic("could not load file")
	}
	query := string(bytes)
	fmt.Println(query)

	cmd := []string{"api", "graphql", "--paginate",
		"-F", "org=" + p.owner,
		"-F", "number=" + p.number,
		"-F", "first=" + "50",
		"-f", "query=" + query,
		"-q", ".data.organization.projectV2"}

	resp := callCLI(cmd)
	fmt.Println(resp)
	return nil
}

func getRepos() {
	bytes, err := GqlFiles.ReadFile("gql/get_repos.gql")
	if err != nil {
		panic("could not load file")
	}
	query := string(bytes)
	fmt.Println(query)

	cmd := []string{"api", "graphql", "--paginate",
		"-f", "query=" + query + ""}

	response := callCLI(cmd)
	fmt.Println(response)
}

func UpdateCreatedAt(p *Project) {
	// fieldId := getFieldId(&p, "Created Date")
	issues := getIssues(p)
	for index, issue := range *issues {
		fmt.Println(index, issue)
	}
}

func generateUpdateStatement(updates []ProjectItem) string {
	t := loadTemplate("gql/update_issues.tmpl")
	var buf bytes.Buffer
	err := t.Execute(&buf, updates)
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
