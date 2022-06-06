package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// New creates a new root command for the CLI.
func New() *cobra.Command {
	c := &cobra.Command{
		Use:           "blockchain-backend",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
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
		},
	}

	c.AddCommand(NewCollector())

	return c
}
