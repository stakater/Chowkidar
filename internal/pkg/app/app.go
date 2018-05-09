package app

import "github.com/stakater/Chowkidar/internal/pkg/cmd"

// Run runs the command
func Run() error {
	cmd := cmd.NewChowkidarCommand()
	return cmd.Execute()
}
