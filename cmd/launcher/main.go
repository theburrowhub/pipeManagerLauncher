// Package main contains the main entrypoint for the application.
package main

import "github.com/sergiotejon/pipeManager/internal/app/launcher/cmd"

// main is the entrypoint for the application
// It sets up the root command and executes the application
func main() {
	cmd.Execute()
}
