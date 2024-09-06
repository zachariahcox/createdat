package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/cli/go-gh/v2"
)

var DEBUG = true
var MAX_UPDATES = 25

//go:embed gql/*
var GqlFiles embed.FS

func main() {
	owner, projectNumber := parseArgs()
	start := time.Now()
	p := NewProject(owner, projectNumber)
	updatedCount := p.UpdateCreatedAt()
	fmt.Println("updated", updatedCount, "items in", time.Since(start))
}

func parseArgs() (string, string) {
	// parse args
	flag.BoolVar(&DEBUG, "debug", true, "Control debug / dry-run mode. No mutations will be made unless this is explicitly set to 'false'.")
	flag.IntVar(&MAX_UPDATES, "maxUpdates", 25, "How many update statements to make in one gql query. Turn this down if you're running into rate limits.")
	url := flag.String("project", "", "fully qualified url to the project")
	flag.Parse()

	if DEBUG {
		fmt.Println("DEBUG:")
		fmt.Println("\tDEBUG =", DEBUG)
		fmt.Println("\tMAX_UPDATES =", MAX_UPDATES)
		fmt.Println("\tURL = ", url)
	}

	// parse url
	if url == nil {
		return "", ""
	}

	components := strings.Split(*url, "/") // https://github.com/user/repo/projects/number

	return "", ""
}

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
	name := path.Base(filePath) // go-gotcha: the template name has to be the _basename_ of "one of the parsed files"
	t, err := template.New(name).ParseFiles(filePath)
	if err != nil {
		log.Fatal("could not load template at path '", filePath, "'", err)
	}
	return t
}
