package main

import (
	"evolved-commit/cmd"
)

func main() {
	// Execute the root command of the CLI application.
	// All errors are handled by Cobra's default error handler, including exiting on failure.
	cmd.Execute()
}
