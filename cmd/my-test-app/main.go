package main

import "log"

// Version injected at compile time.
var version = "No version provided"

func main() {
	log.Printf("Starting my-test-app %s\n", version)
}
