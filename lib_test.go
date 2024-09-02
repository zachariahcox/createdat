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
	result := callCLI(args)
	want := regexp.MustCompile("$gh.*")
	assert(t, want.MatchString(result), strings.Join(args, " "))
}

func TestParseProject(t *testing.T) {
	data := []byte(`{
	"id":"PVT_kwDOBMtIU84AmrDl",
	"title":"SLSA Source Track"
	}`)

	var dat map[string]interface{}
	if err := json.Unmarshal(data, &dat); err != nil {
		t.Fatal(err)
	}
	p := new(Project)
	p.id = dat["id"].(string)
	p.title = dat["title"].(string)
}
func TestParseProjectItem(t *testing.T) {
	jsonData := []byte(`
	{
	"nodes":
	[
		{
		"content": {
			"closed": false,
			"closedAt": null,
			"createdAt": "2024-08-14T18:01:45Z",
			"number": 1111,
			"title": "source track: create a \"levels\" table for the source track",
			"url": "https://github.com/slsa-framework/slsa/issues/1111"
		},
		"createdAt": "2024-08-23T20:49:46Z",
		"id": "PVTI_lADOBMtIU84AmrDlzgSLFHY",
		"type": "ISSUE"
		},
		{
		"content": {
			"closed": true,
			"closedAt": "2024-08-15T15:01:11Z",
			"createdAt": "2024-07-12T18:27:05Z",
			"number": 1097,
			"title": "content: source track draft: simplify and clarify level goals",
			"url": "https://github.com/slsa-framework/slsa/pull/1097"
		},
		"createdAt": "2024-08-23T20:49:46Z",
		"id": "PVTI_lADOBMtIU84AmrDlzgSLFHc",
		"type": "PULL_REQUEST"
		}
	]
	}
	`)
	var dat map[string]interface{}
	if err := json.Unmarshal(jsonData, &dat); err != nil {
		t.Fatal(err)
	}

	nodes := dat["nodes"].([]interface{})
	items := make([]PI, len(nodes), len(nodes))
	for i, node := range nodes {
		itemData := node.(map[string]interface{})
		item := items[i]

		// parse project data
		item.id = itemData["id"].(string)
		item.createdAt = itemData["createdAt"].(string)
		item.contentType = itemData["type"].(string)

		// parse content data
		content := itemData["content"].(map[string]interface{})
		// item.content.id = content["id"].(string)
		item.content.url = content["url"].(string)
		item.content.number = content["number"].(float64)
		item.content.title = content["title"].(string)
		item.content.createdAt = content["createdAt"].(string)
		item.content.closed = content["closed"].(bool)
		if item.content.closed {
			item.content.closedAt = content["closedAt"].(string)
		}
	}
}
