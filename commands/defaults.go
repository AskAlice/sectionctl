package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/alecthomas/kong"
)

// DefaultsCmd handles setting defaults for the command line
type DefaultsCmd struct {
	Set  DefaultsSetCommand  `cmd help:"Set a default"`
	List DefaultsListCommand `cmd help:"List set defaults"`
}

// DefaultsSetCommand records a default
type DefaultsSetCommand struct{}

// Run executes the command
func (c *DefaultsSetCommand) Run() (err error) {
	return err
}

// DefaultsListCommand lists recorded defaults
type DefaultsListCommand struct{}

// Run executes the command
func (c *DefaultsListCommand) Run() (err error) {
	usr, err := user.Current()
	if err != nil {
		return err
	}
	defaultsPath := filepath.Join(usr.HomeDir, ".config", "section", "defaults.json")

	def, err := readDefaults(defaultsPath)
	if err != nil {
		return err
	}

	table := NewTable(os.Stdout)
	table.SetHeader([]string{"Default", "Value"})

	v := reflect.ValueOf(def)
	types := v.Type()
	for i := 0; i < v.NumField(); i++ {
		r := []string{}

		switch t := v.Field(i).Interface().(type) {
		case int:
			if v.Field(i).Interface().(int) == 0 {
				continue
			}
			r = []string{types.Field(i).Name, strconv.Itoa(v.Field(i).Interface().(int))}
		default:
			return fmt.Errorf("unhandled default type %s", t)
		}

		table.Append(r)
	}

	table.Render()
	return err
}

// readDefaults loads and returns sectionctl defaults
func readDefaults(path string) (def defaults, err error) {
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return def, err
	}
	defaultsFile, err := os.Open(path)
	if err != nil {
		return def, err
	}
	defer defaultsFile.Close()

	c, err := ioutil.ReadAll(defaultsFile)
	if err != nil {
		return def, err
	}

	err = json.Unmarshal(c, &def)
	return def, err
}

// defaults are defaults that can be set and queries
type defaults struct {
	AccountID     int `json:"account_id"`
	ApplicationID int `json:"application_id"`
}

// DefaultAccountIDResolver looks up a default Section account id
var DefaultAccountIDResolver kong.ResolverFunc = func(context *kong.Context, parent *kong.Path, flag *kong.Flag) (interface{}, error) {
	if flag.Name == "account-id" {
		usr, err := user.Current()
		if err != nil {
			return nil, nil
		}
		defaultsPath := filepath.Join(usr.HomeDir, ".config", "section", "defaults.json")
		def, err := readDefaults(defaultsPath)
		if err != nil {
			return nil, nil
		}

		if def.AccountID <= 0 {
			return nil, nil
		}
		return def.AccountID, nil
	}
	return nil, nil
}
