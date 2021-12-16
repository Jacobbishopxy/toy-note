package main

import (
	"toy-note/logger"
)

const logPath = "../../../logs/toy-note.log"

func main() {
	// Initialize logger
	if err := logger.Init("debug", logPath, false); err != nil {
		panic(err)
	} else {
		defer logger.TNLogger.Sync()
	}

	// Main code here...
	slog := logger.TNLogger.NewSugar("Main")

	slog.Info("Hello world!")
}
