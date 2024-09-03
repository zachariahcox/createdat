package main

import (
	"flag"
	"fmt"
	"time"
)

var DEBUG = true
var MAX_UPDATES = 25

func main() {
	// parse args
	flag.BoolVar(&DEBUG, "debug", true, "Control debug / dry-run mode. No mutations will be made unless this is explicitly set to 'false'.")
	flag.IntVar(&MAX_UPDATES, "maxUpdates", 25, "How many update statements to make in one gql query. Turn this down if you're running into rate limits.")
	flag.Parse()

	fmt.Println("DEBUG =", DEBUG)
	fmt.Println("MAX_UPDATES =", MAX_UPDATES)

	start := time.Now()

	p := NewProject("slsa-framework", "slsa", "5")
	updatedCount := p.UpdateCreatedAt()

	fmt.Println("updated", updatedCount, "items in", time.Since(start))
}

func track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func duration(msg string, start time.Time) {
	fmt.Printf("%v: %v\n", msg, time.Since(start))
}
