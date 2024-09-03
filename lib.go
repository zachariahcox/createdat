package main

import (
	"embed"
	"fmt"
	"log"
	"path"
	"strings"
	"text/template"

	"github.com/cli/go-gh/v2"
)

//go:embed gql/*
var GqlFiles embed.FS

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
	name := path.Base(filePath) // the template name has to be the basename of "one of the files"
	t, err := template.New(name).ParseFiles(filePath)
	if err != nil {
		panic(err)
	}
	return t
}

func getRepos() {
	bytes, err := GqlFiles.ReadFile("gql/get_repos.gql")
	if err != nil {
		panic("could not load file")
	}
	query := string(bytes)
	fmt.Println(query)

	cmd := []string{"api", "graphql", "--paginate",
		"-f", "query=" + query + ""}

	response := callCLI(cmd)
	fmt.Println(response)
}
