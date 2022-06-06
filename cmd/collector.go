package cmd

import "github.com/spf13/cobra"

func NewCollector() *cobra.Command {
	c := &cobra.Command{
		Use:     "collector [command]",
		Short:   "Commands for managing the collector service",
		Long:    "The collector service saves Tendermint TXs and events in a data backend",
		Aliases: []string{"c"},
		Args:    cobra.ExactArgs(1),
	}

	c.AddCommand(NewCollectorStart())

	return c
}
