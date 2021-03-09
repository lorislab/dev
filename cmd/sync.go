package cmd

import (
	"sync"

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

			count := 0
			sum := 0

			for _, priority := range priorities {
				var wg sync.WaitGroup
				count = 0

				for _, app := range apps[priority] {
					count++
					sum++
					wg.Add(1)
					go env.Sync(app, &wg, flags.ForceUpdate, true)
				}
				wg.Wait()

				log.Info().Int("count", count).Int("sum", sum).Int("priority", priority).Msg("Sync apps finished")
			}
			log.Info().Int("sum", sum).Msg("Sync apps finished.")
		},
	}

	addBoolFlag(cmd, "force-update", "", false, "force update repositories before sync")

	return cmd
}
