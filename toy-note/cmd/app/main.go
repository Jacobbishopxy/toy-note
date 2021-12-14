package main

import "toy-note/logger"

const logPath = "../../../logs/toy-note.log"

func main() {
	// Initialize logger
	if err := logger.Init("debug", logPath, false); err != nil {
		panic(err)
	} else {
		defer logger.Sync()
	}

	// Main code here...
	slog := logger.NewSugar("Main")

	slog.Info("Hello world!")
}
