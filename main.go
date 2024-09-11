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

func main() {
	start := time.Now()

	// parse args
	flag.BoolVar(&DEBUG, "debug", true, "Control debug / dry-run mode. No mutations will be made unless this is explicitly set to 'false'.")
	flag.IntVar(&MAX_UPDATES, "maxUpdates", 25, "How many update statements to make in one gql query. Turn this down if you're running into rate limits.")
	url := flag.String("project", "", "fully qualified url to the project")
	flag.Parse()
	if DEBUG {
		fmt.Println("DEBUG:")
		fmt.Println("\tDEBUG =", DEBUG)
		fmt.Println("\tMAX_UPDATES =", MAX_UPDATES)
		fmt.Println("\tURL = " + *url)
	}
	if url == nil {
		log.Fatal("you must provide a project url")
	}

	// load project
	scope, owner, number := parseUrl(*url)
	p := NewProject(scope, owner, number)

	updatedCount := p.UpdateCreatedAt()
	fmt.Println("updated", updatedCount, "items in", time.Since(start))
}

func parseUrl(url string) (string, string, string) {
	components := strings.Split(url, "/")
	for i, c := range components {
		if c == "github.com" {
			return strings.TrimSuffix(components[i+1], "s"), // the url will have orgs or users
				components[i+2], // owner name
				components[i+4] // project number
		}
	}
	return "", "", ""
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

//go:embed gql/*
var GqlFiles embed.FS

func loadTemplate(filePath string) *template.Template {
	// go-gotcha: the template name has to be the _basename_ of "one of the parsed files"
	name := path.Base(filePath)
	t, err := template.New(name).ParseFS(GqlFiles, filePath)
	if err != nil {
		log.Fatal("could not load template at path '", filePath, "'", err)
	}
	return t
}
