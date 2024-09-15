package main

import (
	"encoding/json"
	"log"
)

type Issue struct {
	Id         string `json:"id,omitempty"`
	Url        string `json:"url,omitempty"`
	Title      string `json:"title,omitempty"`
	CreatedAt  string `json:"createdAt,omitempty"`
	UpdatedAt  string `json:"updatedAt,omitempty"`
	Repository struct {
		Url string
	} `json:"repository,omitempty"`
	Assignees struct {
		Nodes []struct {
			Login string
		} `json:"nodes,omitempty"`
	} `json:"assignees,omitempty"`
	Labels struct {
		Nodes []struct {
			Name string
		} `json:"nodes,omitempty"`
	} `json:"labels,omitempty"`
}

type Repository struct {
	Owner  string
	Name   string
	Issues []Issue `json:"items,omitempty"`
}

func NewRepository(owner, name string) *Repository {
	r := &Repository{
		Owner: owner,
		Name:  name,
	}
	r.LoadIssues()
	return r
}
func (r *Repository) nwo() string {
	return r.Owner + "/" + r.Name
}

func (r *Repository) LoadIssues() {
	query := loadQuery("gql/get_issues.gql")
	filter := "repo:" + r.nwo() + " is:issue is:open"
	cmd := []string{"api", "graphql", "--paginate",
		"-F", "filters=" + filter,
		"-F", "first=" + "50",
		"-f", "query=" + query,
		"-q", ".data.search.nodes"}

	resp := callCLI(cmd)
	if resp == nil {
		log.Fatal("could not load issues")
		return
	}

	json.Unmarshal(resp, &r.Issues)
}
