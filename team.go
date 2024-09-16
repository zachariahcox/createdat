package main

import (
	"encoding/csv"
	"os"
)

type IssueQuery struct {
	OrgName  string
	RepoName string
	Query    string
}
type Team struct {
	Org          string
	Name         string
	ServiceNames []string
	IssueQueries []IssueQuery
}

func ParseTeams(filepath string) ([]*Team, error) {
	// Open CSV file
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all rows from the CSV
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Initialize the slice to hold the teams
	var teams []*Team
	// Iterate over the rows and populate the teams slice
	for _, row := range rows[1:] { // Skip the header row
		org := row[0]
		name := row[1]
		var team *Team
		for _, t := range teams {
			if t.Org == org && t.Name == name {
				team = t
				break
			}
		}
		if team == nil {
			// If the team doesn't exist, create a new one
			teams = append(teams, &Team{
				Org:          org,
				Name:         name,
				ServiceNames: make([]string, 0, 10),
				IssueQueries: make([]IssueQuery, 0, 10),
			})
			team = teams[len(teams)-1]
		}

		//
		team.ServiceNames = append(team.ServiceNames, row[2]) // append the service name to the team
		team.IssueQueries = append(team.IssueQueries, IssueQuery{row[3], row[4], row[5]})
	}

	return teams, nil
}
