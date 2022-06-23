package cmd

import (
	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	"github.com/spf13/cobra"

	"github.com/ignite/backend/collector"
)

const (
	flagCollectTimeout = "collect-timeout"
	flagGracePeriod    = "grace-period"
	flagMinHeight      = "min-height"
	flagRPCAddress     = "rpc-address"

	defaultRPCAddress = "http://127.0.0.1:26657"
)

func NewCollectorStart() *cobra.Command {
	c := &cobra.Command{
		Use:   "start",
		Short: "Start the collector service",
		RunE:  collectorStartHandler,
	}

	c.Flags().AddFlagSet(flagSetDatabase())
	c.Flags().Duration(flagCollectTimeout, collector.DefaultTimeout, "Collect timeout")
	c.Flags().Duration(flagGracePeriod, collector.DefaultGrace, "Grace period between collect calls")
	c.Flags().Int64(flagMinHeight, collector.DefaultMinHeight, "Minimum block height to start collecting from")
	c.Flags().StringP(flagRPCAddress, "a", defaultRPCAddress, "RPC address of the chain")

	c.MarkFlagRequired(flagDatabaseName)

	return c
}

func collectorStartHandler(cmd *cobra.Command, args []string) error {
	rpcAddr, err := cmd.Flags().GetString(flagRPCAddress)
	if err != nil {
		return err
	}

	minHeight, err := cmd.Flags().GetInt64(flagMinHeight)
	if err != nil {
		return err
	}

	grace, err := cmd.Flags().GetDuration(flagGracePeriod)
	if err != nil {
		return err
	}

	timeout, err := cmd.Flags().GetDuration(flagCollectTimeout)
	if err != nil {
		return err
	}

	db, err := createDatabaseAdapter(cmd)
	if err != nil {
		return err
	}

	ctx := cmd.Context()

	client, err := cosmosclient.New(ctx, cosmosclient.WithNodeAddress(rpcAddr))
	if err != nil {
		return err
	}

	service := collector.NewService(
		db,
		client,
		collector.WithTimeout(timeout),
		collector.WithMinHeight(minHeight),
		collector.WithGrace(grace),
	)

	return service.Run(ctx)
}
