package main

import (
	"encoding/json"
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
	result := string(callCLI(args))
	want := regexp.MustCompile("$gh.*")
	assert(t, want.MatchString(result), strings.Join(args, " "))
}

func TestParseProject(t *testing.T) {
	resp := []byte(`
	{
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

	var data = new(Project)
	if err := json.Unmarshal(resp, &data); err != nil {
		t.Fail()
	}
	assert(t, data != nil, "unmarshal should not be nil")
	assert(t, len(data.Items.Nodes) == 2, "should be 2 items")
}
