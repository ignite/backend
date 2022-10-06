package cmd

import (
	"errors"
	"strings"

	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/adapter/postgres"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultLogFormat = logFormatText
	defaultLogLevel  = "info"

	flagDatabaseHost     = "database-host"
	flagDatabaseName     = "database-name"
	flagDatabaseParams   = "database-params"
	flagDatabasePassword = "database-password"
	flagDatabasePort     = "database-port"
	flagDatabaseUser     = "database-user"
	flagLogFormat        = "log-format"
	flagLogLevel         = "log-level"

	logFormatJSON = "json"
	logFormatText = "text"
)

// New creates a new root command for the CLI.
func New() *cobra.Command {
	c := &cobra.Command{
		Use:           "ignite-backend",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			initViper(cmd)

			if err := initLogging(cmd); err != nil {
				return err
			}

			return nil
		},
	}

	c.PersistentFlags().AddFlagSet(flagSetLogging())

	c.AddCommand(NewAPI())
	c.AddCommand(NewCollector())

	return c
}

func initViper(cmd *cobra.Command) {
	// Prepare Viper to read values from the environment
	viper.SetEnvPrefix("IGN")
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

func initLogging(cmd *cobra.Command) error {
	name, err := cmd.Flags().GetString(flagLogLevel)
	if err != nil {
		return err
	}

	level, err := log.ParseLevel(name)
	if err != nil {
		return err
	}

	log.SetLevel(level)

	format, err := cmd.Flags().GetString(flagLogFormat)
	if err != nil {
		return err
	}

	switch format {
	case logFormatJSON:
		log.SetFormatter(&log.JSONFormatter{})
	case logFormatText:
		// Text is the default one for logrus
	default:
		return errors.New("invalid log format")
	}

	return nil
}

func flagSetLogging() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(flagLogFormat, defaultLogFormat, "Log format [json|text]")
	fs.String(flagLogLevel, defaultLogLevel, "Log level [trace|debug|info|warn|error|fatal|panic]")

	return fs
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
