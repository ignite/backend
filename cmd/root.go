package cmd

import (
	"strings"

	"github.com/ignite-hq/cli/ignite/pkg/cosmosmetric/adapter/postgres"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	flagDatabaseHost     = "database-host"
	flagDatabaseName     = "database-name"
	flagDatabaseParams   = "database-params"
	flagDatabasePassword = "database-password"
	flagDatabasePort     = "database-port"
	flagDatabaseUser     = "database-user"
)

// New creates a new root command for the CLI.
func New() *cobra.Command {
	c := &cobra.Command{
		Use:           "ignite-backend",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initViper(cmd)
		},
	}

	c.AddCommand(NewAPI())
	c.AddCommand(NewCollector())

	return c
}

func initViper(cmd *cobra.Command) {
	// Prepare Viper to read values from the environment
	viper.SetEnvPrefix("BCB")
	viper.AutomaticEnv()

	// Initialize the flags with the values setted in the environment.
	// This is done to be able to mark flags as required allowing their
	// initialization to be done thought the environment.
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Name == "help" {
			return
		}

		// Viper environment variable name for the flag
		name := strings.ReplaceAll(f.Name, "-", "_")

		// Assign the env value to the flag
		if viper.IsSet(name) && viper.GetString(name) != "" {
			cmd.Flags().Set(f.Name, viper.GetString(name))
		}
	})
}

func flagSetDatabase() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.StringP(flagDatabaseHost, "H", postgres.DefaultHost, "Database server hostname or IP")
	fs.StringP(flagDatabaseName, "d", "", "Name of the database where to store the data")
	fs.StringSliceP(flagDatabaseParams, "P", nil, "Extra database parameters [name=value,...]")
	fs.String(flagDatabasePassword, "", "Database user password")
	fs.UintP(flagDatabasePort, "p", postgres.DefaultPort, "Database server port")
	fs.StringP(flagDatabaseUser, "U", "", "Database user name")

	return fs
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

func createDatabaseAdapter(cmd *cobra.Command) (postgres.Adapter, error) {
	database, err := cmd.Flags().GetString(flagDatabaseName)
	if err != nil {
		return postgres.Adapter{}, err
	}

	host, err := cmd.Flags().GetString(flagDatabaseHost)
	if err != nil {
		return postgres.Adapter{}, err
	}

	port, err := cmd.Flags().GetUint(flagDatabasePort)
	if err != nil {
		return postgres.Adapter{}, err
	}

	user, err := cmd.Flags().GetString(flagDatabaseUser)
	if err != nil {
		return postgres.Adapter{}, err
	}

	password, err := cmd.Flags().GetString(flagDatabasePassword)
	if err != nil {
		return postgres.Adapter{}, err
	}

	params, err := cmd.Flags().GetStringSlice(flagDatabaseParams)
	if err != nil {
		return postgres.Adapter{}, err
	}

	return postgres.NewAdapter(
		database,
		postgres.WithHost(host),
		postgres.WithPort(port),
		postgres.WithUser(user),
		postgres.WithPassword(password),
		postgres.WithParams(parseDatabaseParamsFlag(params)),
	)
}
