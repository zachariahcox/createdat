package main

func main() {
	p := Project{"slsa-framework", "slsa", "5"}
	// UpdateCreatedAt()
	getIssuesMissingDates(&p)
	// getRepos()
}
