package cmd

import (
	"github.com/spf13/cobra"
)

type command interface {
	execute() error
}

type cobraRunSignature func(cmd *cobra.Command, args []string) error

type baseCommand struct {
	c          conf
	path       string
	top        int
	removeRoot bool
}

func newBaseCmd(c conf, path string) baseCommand {
	return baseCommand{
		c:          c,
		top:        *c.globals().top,
		removeRoot: *c.globals().removeRoot,
		path:       path,
	}
}

type cobraCreator struct {
	createCmd func() command
}

func (c *cobraCreator) runE() cobraRunSignature {
	return func(cmd *cobra.Command, args []string) error {
		return c.createCmd().execute()
	}
}

func (c *cobraCreator) newCobraCommand(use, alias, short string) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     use,
		Aliases: []string{alias},
		Short:   short,
		RunE:    c.runE(),
	}
	return cmd
}
