package cmd

import (
	"dirstat/module"
	"github.com/spf13/cobra"
)

func newBenford(c conf) *cobra.Command {
	opt := options{}

	var cmd = &cobra.Command{
		Use:     "b",
		Aliases: []string{"benford"},
		Short:   "Show the first digit distribution of files size (benford law validation)",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := module.NewContext(*c.globals().top, *c.globals().removeRoot, opt.path)
			benford := module.NewBenfordFileModule(ctx)
			totalmod := module.NewTotalModule(ctx)

			run(opt.path, c, benford, totalmod)

			return nil
		},
	}

	configure(cmd, &opt)

	return cmd
}
