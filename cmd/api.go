package cmd

import "github.com/spf13/cobra"

func NewAPI() *cobra.Command {
	c := &cobra.Command{
		Use:     "api [command]",
		Short:   "Commands for managing the API service",
		Long:    "The API service exposes an RPC API to query events and TXs",
		Aliases: []string{"c"},
		Args:    cobra.ExactArgs(1),
	}

	c.AddCommand(NewAPIStart())

	return c
}
