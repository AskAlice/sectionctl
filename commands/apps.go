package commands

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/section/sectionctl/api"
)

// AppsCmd manages apps on Section
type AppsCmd struct {
	List   AppsListCmd   `cmd help:"List apps on Section." default:"1"`
	Info   AppsInfoCmd   `cmd help:"Show detailed app information on Section."`
	Create AppsCreateCmd `cmd help:"Create new app on Section."`
	Delete AppsDeleteCmd `cmd help:"Delete an existing app on Section."`
}

// AppsListCmd handles listing apps running on Section
type AppsListCmd struct {
	AccountID int `required short:"a"`
}

// NewTable returns a table with sectionctl standard formatting
func NewTable(out io.Writer) (t *tablewriter.Table) {
	t = tablewriter.NewWriter(out)
	t.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	t.SetCenterSeparator("|")
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	return t
}

// Run executes the command
func (c *AppsListCmd) Run() (err error) {
	s := NewSpinner("Looking up apps")
	s.Start()

	apps, err := api.Applications(c.AccountID)
	s.Stop()
	if err != nil {
		return err
	}

	table := NewTable(os.Stdout)
	table.SetHeader([]string{"App ID", "App Name"})

	for _, a := range apps {
		r := []string{strconv.Itoa(a.ID), a.ApplicationName}
		table.Append(r)
	}

	table.Render()
	return err
}

// AppsInfoCmd shows detailed information on an app running on Section
type AppsInfoCmd struct {
	AccountID int `required short:"a"`
	AppID     int `required short:"i"`
}

// Run executes the command
func (c *AppsInfoCmd) Run() (err error) {
	s := NewSpinner("Looking up app info")
	s.Start()

	app, err := api.Application(c.AccountID, c.AppID)
	s.Stop()
	if err != nil {
		return err
	}

	fmt.Printf("🌎🌏🌍\n")
	fmt.Printf("App Name: %s\n", app.ApplicationName)
	fmt.Printf("App ID: %d\n", app.ID)
	fmt.Printf("Environment count: %d\n", len(app.Environments))

	for i, env := range app.Environments {
		fmt.Printf("\n-----------------\n\n")
		fmt.Printf("Environment #%d: %s (ID:%d)\n\n", i+1, env.EnvironmentName, env.ID)
		fmt.Printf("💬 Domains (%d total)\n", len(env.Domains))

		for _, dom := range env.Domains {
			fmt.Println()

			table := NewTable(os.Stdout)
			table.SetHeader([]string{"Attribute", "Value"})
			table.SetAutoMergeCells(true)
			r := [][]string{
				[]string{"Domain name", dom.Name},
				[]string{"Zone name", dom.ZoneName},
				[]string{"CNAME", dom.CNAME},
				[]string{"Mode", dom.Mode},
			}
			table.AppendBulk(r)
			table.Render()
		}

		fmt.Println()
		mod := "modules"
		if len(env.Stack) == 1 {
			mod = "module"
		}
		fmt.Printf("🥞 Stack (%d %s total)\n", len(env.Stack), mod)
		fmt.Println()

		table := NewTable(os.Stdout)
		table.SetHeader([]string{"Name", "Image"})
		table.SetAutoMergeCells(true)
		for _, p := range env.Stack {
			r := []string{p.Name, p.Image}
			table.Append(r)
		}
		table.Render()
	}

	fmt.Println()

	return err
}

// AppsCreateCmd handles creating apps on Section
type AppsCreateCmd struct {
	AccountID int    `required short:"a" help:"ID of account to create the app under"`
	Hostname  string `required short:"d" help:"FQDN the app can be accessed at"`
	Origin    string `required short:"o" help:"URL to fetch the origin"`
	StackName string `required short:"s" help:"Name of stack to deploy"`
}

// Run executes the command
func (c *AppsCreateCmd) Run() (err error) {
	s := NewSpinner(fmt.Sprintf("Creating new app %s", c.Hostname))
	s.Start()

	api.Timeout = 120 * time.Second // this specific request can take a long time
	r, err := api.ApplicationCreate(c.AccountID, c.Hostname, c.Origin, c.StackName)
	s.Stop()
	if err != nil {
		return err
	}

	fmt.Printf("\nSuccess: created app '%s' with id '%d'\n", r.ApplicationName, r.ID)

	return err
}

// AppsDeleteCmd handles deleting apps on Section
type AppsDeleteCmd struct {
	AccountID int `required short:"a" help:"ID of account the app belongs to"`
	AppID     int `required short:"i" help:"ID of the app to delete"`
}

// Run executes the command
func (c *AppsDeleteCmd) Run() (err error) {
	s := NewSpinner(fmt.Sprintf("Deleting app with id '%d'", c.AppID))
	s.Start()

	api.Timeout = 120 * time.Second // this specific request can take a long time
	_, err = api.ApplicationDelete(c.AccountID, c.AppID)
	s.Stop()
	if err != nil {
		return err
	}

	fmt.Printf("\nSuccess: deleted app with id '%d'\n", c.AppID)

	return err
}
