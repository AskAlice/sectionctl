package commands

import (
	"fmt"

	"github.com/section/sectionctl/api"
	"github.com/section/sectionctl/api/auth"
)

// LogoutCmd handles revoking previously set up authentication
type LogoutCmd struct{}

// Run executes the command
func (c *LogoutCmd) Run() (err error) {
	s := NewSpinner(fmt.Sprintf("Revoking your authentication for %s", api.PrefixURI.Host))
	s.Start()
	err = auth.DeleteCredential(api.PrefixURI.Host)
	s.Stop()
	return err
}
