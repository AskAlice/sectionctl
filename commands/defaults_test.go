package commands

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandsDefaultsReadDefaultsReturnsDefaults(t *testing.T) {
	assert := assert.New(t)

	path := filepath.Join("testdata", "defaults", "defaults.valid.json")

	// Invoke
	d, err := readDefaults(path)

	// Test
	assert.NoError(err)
	assert.Greater(d.AccountID, 0)
	assert.Greater(d.ApplicationID, 0)
}

func TestCommandsDefaultsDefaultAccountIDResolver(t *testing.T) {
	assert := assert.New(t)
	assert.FailNow("no")
}
