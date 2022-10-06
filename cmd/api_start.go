package cmd

import (
	"github.com/ignite/backend/api"
	"github.com/spf13/cobra"
)

const (
	flagAddress = "address"
)

func NewAPIStart() *cobra.Command {
	c := &cobra.Command{
		Use:   "start",
		Short: "Start the API service",
		RunE:  apiStartHandler,
	}

	c.Flags().AddFlagSet(flagSetDatabase())
	c.Flags().StringP(flagAddress, "a", api.DefaultAddress, "address to listen for requests")
	c.MarkFlagRequired(flagDatabaseName)

	return c
}

func apiStartHandler(cmd *cobra.Command, args []string) error {
	addr, err := cmd.Flags().GetString(flagAddress)
	if err != nil {
		return err
	}

	db, err := createDatabaseAdapter(cmd)
	if err != nil {
		return err
	}

	service := api.NewService(db, api.WithAddress(addr))

	return service.Run(cmd.Context())
}
