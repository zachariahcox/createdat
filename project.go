package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// kind of a cool facility here: https://mholt.github.io/json-to-go/

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

func NewProject(owner string, repo string, number string) *Project {
	p := new(Project)
	p.Owner = owner
	p.Repo = repo
	p.Number = number
	p.UpdateItems()
	return p
}

func (p *Project) UpdateItems() {
	b, err := GqlFiles.ReadFile("gql/get_project_contents.gql")
	if err != nil {
		panic("could not load file")
	}
	query := string(b)
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

func (p *Project) UpdateFields() {
	b, err := GqlFiles.ReadFile("gql/get_project_fields.gql")
	if err != nil {
		panic("could not load file")
	}
	query := string(b)
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

func (p *Project) GetFieldId(fieldName string) (int, string) {
	cmd := []string{"project", "field-list",
		"--owner", p.Owner, p.Number,
		"--format", "json",
		"--jq", ".fields[] | select(.name==\"" + fieldName + "\") | .id"}
	fieldId := callCLI(cmd)
	if fieldId == nil {
		log.Fatal("Could not find field named '" + fieldName + "'")
	}

	return 0, strings.TrimSuffix(string(fieldId), "\n")
}

type ProjectItemUpdate struct {
	FieldIndex          int
	ProjectIndex        int
	ProjectId           string
	ProjectItemId       string
	FieldId             string
	ProjectV2FieldValue string // this is an https://docs.github.com/en/graphql/reference/input-objects#projectv2fieldvalue
}

func generateUpdateStatement(updates []*ProjectItemUpdate) string {
	t := loadTemplate("gql/update_issues.tmpl")
	var buf bytes.Buffer
	err := t.Execute(&buf, updates)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func (p *Project) UpdateCreatedAt() int {
	fieldIndex, fieldId := p.GetFieldId("Created Date")
	updates := make([]*ProjectItemUpdate, 0, len(p.Items.Nodes))
	for itemIndex, item := range p.Items.Nodes {
		// does this item have a created date?
		hasCreatedAt := false
		for _, f := range item.FieldValues.Nodes {
			hasCreatedAt = f.Field.ID == fieldId
			if hasCreatedAt {
				break
			}
		}

		if !hasCreatedAt {
			update := new(ProjectItemUpdate)
			update.FieldId = fieldId
			update.FieldIndex = fieldIndex
			update.ProjectId = p.ID
			update.ProjectIndex = itemIndex
			update.ProjectItemId = item.ID
			update.ProjectV2FieldValue = "date:\"" + item.Content.CreatedAt + "\"" // crazy memex syntax
			updates = append(updates, update)
		}
	}

	len_updates := len(updates)
	for i := 0; i < len_updates; i += MAX_UPDATES {

		// golang HAS NO MIN FUNCTION FOR INTEGERS.
		end := i + MAX_UPDATES
		if end > len_updates {
			end = len_updates
		}

		s := generateUpdateStatement(updates[i:end])
		cmd := []string{"api", "graphql", "-f", "query=" + s}
		if DEBUG {
			fmt.Println("gh", cmd)
		} else {
			callCLI(cmd)
		}
	}

	return len_updates
}
