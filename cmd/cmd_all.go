package cmd

import (
	"dirstat/module"
	"github.com/spf13/cobra"
)

func newAll(c conf) *cobra.Command {
	opt := options{}

	var cmd = &cobra.Command{
		Use:     "a",
		Aliases: []string{"all"},
		Short:   "Show all information about folder/volume",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := module.NewContext(*c.globals().top, *c.globals().removeRoot, opt.path)
			foldersmod := module.NewFoldersModule(ctx)
			totalmod := module.NewTotalModule(ctx)
			detailfilemod := module.NewDetailFileModule(ctx, opt.vrange)
			totalfilemod := module.NewAggregateFileModule(ctx)
			extmod := module.NewExtensionModule(ctx)
			topfilesmod := module.NewTopFilesModule(ctx)
			benford := module.NewBenfordFileModule(ctx)

			run(opt.path, c, totalfilemod, benford, extmod, topfilesmod, foldersmod, detailfilemod, totalmod)

			return nil
		},
	}

	configure(cmd, &opt)

	return cmd
}
