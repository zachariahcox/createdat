package main

import (
	"bytes"
	"embed"
	"encoding/json"
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

func callCLI(cmd []string) []byte {
	stdout, stderr, err := gh.Exec(cmd...)
	if err != nil {
		log.Fatal(strings.Join(cmd, " "), "\n",
			stdout.String(), "\n",
			stderr.String(), "\n",
			err)
		return nil
	}
	return stdout.Bytes()
}

func loadTemplate(filePath string) *template.Template {
	templateFile := filePath
	t, err := template.New(templateFile).ParseFiles(templateFile)
	if err != nil {
		panic(err)
	}
	return t
}

// kind of a cool facility here: https://mholt.github.io/json-to-go/
type ProjectItemGql struct {
	Content struct {
		Closed    bool   `json:"closed,omitempty"`
		ClosedAt  any    `json:"closedAt,omitempty"`
		CreatedAt string `json:"createdAt,omitempty"`
		Number    int    `json:"number,omitempty"`
		Title     string `json:"title,omitempty"`
		URL       string `json:"url,omitempty"`
	} `json:"content,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
	FieldValues struct {
		Nodes []struct {
			Labels struct {
				Nodes []struct {
					Name string `json:"name,omitempty"`
				} `json:"nodes,omitempty"`
			} `json:"labels,omitempty"`
			Field struct {
				ID string `json:"id,omitempty"`
			} `json:"field,omitempty"`
			ID       string `json:"id,omitempty"`
			OptionID string `json:"optionId,omitempty"`
		} `json:"nodes,omitempty"`
	} `json:"fieldValues,omitempty"`
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

type Project struct {
	Owner  string
	Repo   string
	Number string
	ID     string `json:"id,omitempty"`
	Title  string `json:"title,omitempty"`
	Items  struct {
		Nodes []ProjectItemGql `json:"nodes,omitempty"`
		// PageInfo struct {
		// 	EndCursor   string `json:"endCursor,omitempty"`
		// 	HasNextPage bool   `json:"hasNextPage,omitempty"`
		// } `json:"pageInfo,omitempty"`
		TotalCount int `json:"totalCount,omitempty"`
	} `json:"items,omitempty"`
}

func NewProject(owner string, repo string, number string) *Project {
	p := new(Project)
	p.Owner = owner
	p.Repo = repo
	p.Number = number
	getProjectContents(p)
	return p
}

type ProjectItemUpdate struct {
	FieldIndex          int
	ProjectIndex        int
	ProjectId           string
	ProjectItemId       string
	FieldId             string
	ProjectV2FieldValue string // this is an https://docs.github.com/en/graphql/reference/input-objects#projectv2fieldvalue
}

func getFieldId(p *Project, fieldName string) string {
	cmd := []string{"project", "field-list",
		"--owner", p.Owner, p.Number,
		"--format", "json",
		"--jq", ".fields[] | select(.name==\"" + fieldName + "\") | .id"}
	fieldId := callCLI(cmd)
	if fieldId == nil {
		log.Fatal("Could not find field named '" + fieldName + "'")
	}
	return strings.TrimSuffix(string(fieldId), "\n")
}

func getProjectContents(p *Project) {
	bytes, err := GqlFiles.ReadFile("gql/get_project_contents.gql")
	if err != nil {
		panic("could not load file")
	}
	query := string(bytes)
	cmd := []string{"api", "graphql", "--paginate",
		"-F", "org=" + p.Owner,
		"-F", "number=" + p.Number,
		"-F", "first=" + "50",
		"-f", "query=" + query,
		"-q", ".data.organization.projectV2"}

	resp := callCLI(cmd)
	if resp == nil {
		return
	}

	if err := json.Unmarshal(resp, &p); err != nil {
		return
	}
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
	fieldId := getFieldId(p, "Created Date")
	needsUpdate := make([]*ProjectItemUpdate, 0, len(p.Items.Nodes))
	for _, item := range p.Items.Nodes {
		// does this item have a created date?
		hasCreatedAt := false
		for _, f := range item.FieldValues.Nodes {
			hasCreatedAt = f.Field.ID == fieldId
			if hasCreatedAt {
				fmt.Println(item.Content.URL, "already has a createdAt")
				break
			}
		}
		if !hasCreatedAt {
			update := new(ProjectItemUpdate)
			update.FieldId = fieldId
			update.FieldIndex = 0
			needsUpdate = append(needsUpdate, update)
		}
	}

	// generate update statement
	for _, i := range needsUpdate {
		fmt.Println("will update", i.Content.URL, "to", i.Content.CreatedAt)
	}
}

func generateUpdateStatement(updates []ProjectItemUpdate) string {
	t := loadTemplate("gql/update_issues.tmpl")
	var buf bytes.Buffer
	err := t.Execute(&buf, updates)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func makeChanges(updates []ProjectItemUpdate) {
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
