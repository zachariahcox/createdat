package main

import (
	"os"
	"reflect"
	"testing"
)

func TestParseTeams(t *testing.T) {
	// Create a temporary CSV file
	file, err := os.CreateTemp("", "teams_test.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	// Write test data to the CSV file
	csvContent := `Org,Name,ServiceName,OrgName,RepoName,Query
Org1,Team1,Service1,Org1,Repo1,Query1
Org1,Team1,Service2,Org1,Repo2,Query2
Org2,Team2,Service3,Org2,Repo3,Query3
`
	if _, err := file.WriteString(csvContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	file.Close()

	// Call the function to test
	teams, err := ParseTeams(file.Name())
	if err != nil {
		t.Fatalf("ParseTeams returned an error: %v", err)
	}

	// Define the expected result
	expected := []Team{
		{
			Org:          "Org1",
			Name:         "Team1",
			ServiceNames: []string{"Service1", "Service2"},
			IssueQueries: []IssueQuery{
				{OrgName: "Org1", RepoName: "Repo1", Query: "Query1"},
				{OrgName: "Org1", RepoName: "Repo2", Query: "Query2"},
			},
		},
		{
			Org:          "Org2",
			Name:         "Team2",
			ServiceNames: []string{"Service3"},
			IssueQueries: []IssueQuery{
				{OrgName: "Org2", RepoName: "Repo3", Query: "Query3"},
			},
		},
	}

	// Compare the result with the expected result
	if len(teams) != len(expected) {
		t.Errorf("ParseTeams returned %d teams, expected %d", len(teams), len(expected))
	}
	for i := range teams {
		if !reflect.DeepEqual(teams[i], &expected[i]) {
			t.Errorf("ParseTeams returned %+v, expected %+v", teams[i], &expected[i])
		}
	}
}
