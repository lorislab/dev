package cmd

import (
	"sync"

	"github.com/lorislab/dev/env"
	"github.com/lorislab/dev/pkg/api"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

//EnvUninstall environment uninstall command
func EnvUninstall() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall environment",
		Long:  `Uninstall existing environment base from the cluster.`,
		Run: func(cmd *cobra.Command, args []string) {
			flags := readEnvFlags()
			envConfig := envConfig(flags.File)

			apps, priorities := env.LoadApps(envConfig, flags.Tags, flags.Apps, flags.Priorities)

			count := 0
			sum := 0

			for _, priority := range priorities {
				var wg sync.WaitGroup
				count = 0

				for _, app := range apps[priority] {
					count++
					sum++
					wg.Add(1)
					app.Action = api.AppActionUninstall
					go env.Uninstall(app, &wg, true, false)
				}
				wg.Wait()

				log.Info().Int("count", count).Int("sum", sum).Int("priority", priority).Msg("Uninstall apps finished")
			}
			log.Info().Int("sum", sum).Msg("Uninstall apps finished.")
		},
	}

	return cmd
}
