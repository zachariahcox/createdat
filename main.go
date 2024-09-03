package main

var DEBUG = false
var MAX_UPDATES = 25

func main() {
	p := NewProject("slsa-framework", "slsa", "5")
	p.UpdateCreatedAt()
}
