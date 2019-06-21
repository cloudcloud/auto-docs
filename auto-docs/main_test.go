package main

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// CobraCommand is a short-hand type for short-cutting
// work involving commands executed by cobra.
type CobraCommand func(*cobra.Command, []string)

func TestBuildServerCommand(t *testing.T) {
	assert := assert.New(t)

	c := buildServerCommand()
	assert.Equal("server", c.Use, "Server command should be called server.")
	assert.NotEmpty(c.Short, "Short description of Server command should be non-empty.")

	f := c.Run
	_ = f
}

func TestBuildVersionCommand(t *testing.T) {
	assert := assert.New(t)

	c := buildVersionCommand()
	assert.Equal("version", c.Use, "Version command should be called version.")
	assert.NotEmpty(c.Short, "Short description of Version command should be non-empty.")

	f := c.Run
	_ = f
}
