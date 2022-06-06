package cmd

import (
	"strings"

	"github.com/ignite-hq/blockchain-backend/collector"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosclient"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosmetric/adapter/postgres"
	"github.com/spf13/cobra"
)

const (
	flagCollectTimeout   = "collect-timeout"
	flagDatabaseHost     = "database-host"
	flagDatabaseName     = "database-name"
	flagDatabaseParams   = "database-params"
	flagDatabasePassword = "database-password"
	flagDatabasePort     = "database-port"
	flagDatabaseUser     = "database-user"
	flagGracePeriod      = "grace-period"
	flagMinHeight        = "min-height"
	flagRPCAddress       = "rpc-address"

	defaultRPCAddress = "http://127.0.0.1:26657"
)

func NewCollectorStart() *cobra.Command {
	c := &cobra.Command{
		Use:   "start",
		Short: "Start the collector service",
		RunE:  collectorStartHandler,
	}

	c.Flags().Duration(flagCollectTimeout, collector.DefaultTimeout, "Collect timeout")
	c.Flags().StringP(flagDatabaseHost, "H", postgres.DefaultHost, "Database server hostname or IP")
	c.Flags().StringP(flagDatabaseName, "d", "", "Name of the database where to store the data")
	c.Flags().StringSliceP(flagDatabaseParams, "P", nil, "Extra database parameters (name=value,...)")
	c.Flags().String(flagDatabasePassword, "", "Database user password")
	c.Flags().UintP(flagDatabasePort, "p", postgres.DefaultPort, "Database server port")
	c.Flags().StringP(flagDatabaseUser, "U", "", "Database user name")
	c.Flags().Duration(flagGracePeriod, collector.DefaultGrace, "Grace period between collect calls")
	c.Flags().Int64(flagMinHeight, collector.DefaultMinHeight, "Minimum block height to start collecting from")
	c.Flags().StringP(flagRPCAddress, "a", defaultRPCAddress, "RPC address of the chain")

	c.MarkFlagRequired(flagDatabaseName)

	return c
}

func collectorStartHandler(cmd *cobra.Command, args []string) error {
	dbName, err := cmd.Flags().GetString(flagDatabaseName)
	if err != nil {
		return err
	}

	dbHost, err := cmd.Flags().GetString(flagDatabaseHost)
	if err != nil {
		return err
	}

	dbPort, err := cmd.Flags().GetUint(flagDatabasePort)
	if err != nil {
		return err
	}

	dbUser, err := cmd.Flags().GetString(flagDatabaseUser)
	if err != nil {
		return err
	}

	dbPassword, err := cmd.Flags().GetString(flagDatabasePassword)
	if err != nil {
		return err
	}

	dbParams, err := cmd.Flags().GetStringSlice(flagDatabaseParams)
	if err != nil {
		return err
	}

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

	ctx := cmd.Context()

	client, err := cosmosclient.New(ctx, cosmosclient.WithNodeAddress(rpcAddr))
	if err != nil {
		return err
	}

	params := parseDatabaseParamsFlag(dbParams)

	db, err := postgres.NewAdapter(
		dbName,
		postgres.WithHost(dbHost),
		postgres.WithPort(dbPort),
		postgres.WithUser(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.WithParams(params),
	)
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

func parseDatabaseParamsFlag(params []string) map[string]string {
	m := map[string]string{}

	for _, p := range params {
		v := strings.SplitN(p, "=", 2)
		name := v[0]
		value := v[1]

		m[name] = value
	}

	return m
}
