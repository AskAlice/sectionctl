package commands

import (
	"fmt"
)

// LogsCmd tails an application's logs on Section
type LogsCmd struct{}

// Run executes the `login` command
func (a *LogsCmd) Run() (err error) {
	fmt.Printf("🚫 Not implemented")
	return err
}
