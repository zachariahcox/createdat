package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/cli/go-gh/v2"
)

var DEBUG = true
var JUST_GH_CMD = false
var MAX_UPDATES = 25

func main() {
	start := time.Now()

	// parse args
	flag.BoolVar(&DEBUG, "debug", true, "Control debug / dry-run mode. No mutations will be made unless this is explicitly set to 'false'.")
	flag.IntVar(&MAX_UPDATES, "maxUpdates", 25, "How many update statements to make in one gql query. Turn this down if you're running into rate limits.")
	flag.BoolVar(&JUST_GH_CMD, "cli", false, "just print the cli command that would be run")
	url := flag.String("project", "", "fully qualified url to the project")
	flag.Parse()
	if DEBUG {
		fmt.Println("DEBUG:")
		fmt.Println("\tDEBUG =", DEBUG)
		fmt.Println("\tMAX_UPDATES =", MAX_UPDATES)
		fmt.Println("\tURL = " + *url)
	}
	if url == nil || *url == "" {
		flag.Usage()
		os.Exit(1)
	}

	// load project
	scope, owner, number := parseUrl(*url)
	p := NewProject(scope, owner, number)

	updatedCount := p.UpdateCreatedAt()
	if DEBUG {
		fmt.Println("would have updated", updatedCount, "items in", time.Since(start))
	} else {
		fmt.Println("updated", updatedCount, "items in", time.Since(start))
	}
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

func get_debug_cli_command(cmd []string) string {
	copy := make([]string, 0, len(cmd))
	for i, c := range cmd {
		if strings.Contains(c, "query=") {
			copy = append(copy, "query='"+cmd[i][6:]+"'")
		} else {
			copy = append(copy, cmd[i])
		}
	}
	return "gh " + strings.Join(copy, " ")
}

func callCLI(cmd []string) []byte {
	if JUST_GH_CMD {
		fmt.Print(get_debug_cli_command(cmd), "\n\n")
		return nil
	}
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

// NB: magic comment to embed a folder
//
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

func loadQuery(name string) string {
	b, err := GqlFiles.ReadFile(name)
	if err != nil {
		log.Fatal("could not load file", err)
	}

	if DEBUG || JUST_GH_CMD {
		// remove comments
		lines := strings.Split(string(b), "\n")
		no_comments := make([]string, 0, len(lines))
		for _, line := range lines {
			if strings.Contains(line, "#") {
				continue // skip comments
			}
			no_comments = append(no_comments, strings.TrimSpace(line))
		}

		return strings.Join(no_comments, " ")
	} else {
		// it's fine to leave comments in normally
		return string(b)
	}
}
