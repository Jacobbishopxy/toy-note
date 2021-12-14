package main

import "toy-note/logger"

func main() {
	// Initialize logger
	if err := logger.Init("debug", "./logs/toy-note.log", false); err != nil {
		panic(err)
	} else {
		defer logger.Sync()
	}

	// Main code here...
}
