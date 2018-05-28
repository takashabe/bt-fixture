package main

import (
	fixture "github.com/takashabe/bt-fixture"
)

func main() {
	// omit the error handling
	f, _ := fixture.NewFixture("test-project", "test-instance")
	f.Load("testdata/fixture.yaml")
}
