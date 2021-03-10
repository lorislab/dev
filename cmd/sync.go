package cmd

import (
	"github.com/lorislab/dev/env"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type syncFlags struct {
	Env         envFlags `mapstructure:",squash"`
	ForceUpdate bool     `mapstructure:"force-update"`
}

//EnvSync create environment sync command
func EnvSync() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync environment",
		Long:  `Sync existing environment base on the environment configuration.`,
		Run: func(cmd *cobra.Command, args []string) {

			flags := syncFlags{}
			readOptions(&flags)

			envConfig := envConfig(flags.Env.File)

			apps, priorities := env.LoadApps(envConfig, flags.Env.Tags, flags.Env.Apps, flags.Env.Priorities)

			log.Info().Int("priorities", len(priorities)).Msg("Synchronize all applications started.")
			sum, err := env.SyncApps(envConfig, apps, priorities, flags.ForceUpdate)
			if err != nil {
				log.Error().Err(err).Msg("Error synchronize applications!")
			}
			log.Info().Int("sum", sum).Msg("Synchronize all applications finished.")
		},
	}

	addBoolFlag(cmd, "force-update", "", false, "force update repositories before sync")

	return cmd
}
