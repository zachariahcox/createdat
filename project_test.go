package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func assert(t *testing.T, should_be_true bool, message string) {
	if !should_be_true {
		t.Fatal(message)
	}
}

func TestGH(t *testing.T) {
	args := []string{"version"}
	result := string(callCLI(args))
	assert(t, strings.HasPrefix(result, "gh"), "could not find gh tool")
}

func TestGql(t *testing.T) {
	bytes, err := GqlFiles.ReadFile("gql/get_repos.gql")
	assert(t, err == nil, "could not load file from embedded file system")

	query := string(bytes)
	expected := `query($endCursor: String) {
    viewer {
      repositories(first: 100, after: $endCursor) {
        nodes { nameWithOwner }
        pageInfo {
          hasNextPage
          endCursor
        }
      }
    }
  }`
	assert(t, query == expected, "query should match")

	cmd := []string{"api", "graphql", "--paginate",
		"-f", "query=" + query + ""}

	response := callCLI(cmd)
	assert(t, response != nil, "response should not be nil")
}

func TestParseProject(t *testing.T) {
	var data = new(Project)
	if err := json.Unmarshal(GoodProjectResponse, &data); err != nil {
		t.Fail()
	}
	assert(t, data != nil, "unmarshal should not be nil")
	assert(t, len(data.Items.Nodes) == 2, "should be 2 items")
}

func TestParseNil(t *testing.T) {
	p := new(Project)
	err := json.Unmarshal(nil, &p)
	_, ok := err.(*json.SyntaxError)
	assert(t, ok, "should be a syntax error")
	assert(t, err != nil, "this shouldn't work")
	assert(t, len(p.Items.Nodes) == 0, "no nodes!")
}

func TestParseMalformed(t *testing.T) {
	p := new(Project)
	err := json.Unmarshal([]byte("adkjaldkjafahadflkjdaf}}"), &p)
	_, ok := err.(*json.SyntaxError)
	assert(t, ok, "should be a syntax error")
	assert(t, err != nil, "this shouldn't work")
	assert(t, len(p.Items.Nodes) == 0, "no nodes!")
}

func TestAddIssues(t *testing.T) {

	p := NewProject("user", "zachariahcox", "1")
	r := NewRepository("zachariahcox", "test")
	p.AddIssues(r.Issues)
	// assert(t, len(p.Items.Nodes) > 0, "should have more than 0 nodes")
}

var GoodProjectResponse = []byte(`{
	"id":"abc123",
	"title":"my title",
	"items": {
		"pageInfo": {
			"endCursor": "END",
			"hasNextPage": false
			},
		"totalCount": 2,
		"nodes":[
			{
			"content": {
				"closed": false,
				"closedAt": null,
				"createdAt": "2024-08-14T18:01:45Z",
				"number": 1,
				"title": "issue number 1!",
				"url": "https://github.com/test/issues/1"
				},
			"createdAt": "2024-08-23T20:49:46Z",
			"id": "PVTI_abc123",
			"type": "ISSUE"
			},
			{
			"content": {
				"closed": true,
				"closedAt": "2024-08-15T15:01:11Z",
				"createdAt": "2024-07-12T18:27:05Z",
				"number": 2,
				"title": "issue number 2",
				"url": "https://github.com/test/issues/2"
				},
			"createdAt": "2024-08-23T20:49:46Z",
			"id": "PVTI_def345",
			"type": "PULL_REQUEST"
			}
		]
	}
	}`)
