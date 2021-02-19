package cmd

import (
	"fmt"

	"github.com/gosuri/uitable"
	"github.com/lorislab/dev/env"
	"github.com/spf13/cobra"
)

//EnvStatus environment status command
func EnvStatus() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Status of the environment",
		Long:  `Status of the environment.`,
		Run: func(cmd *cobra.Command, args []string) {

			flags := readEnvFlags()
			envConfig := envConfig(flags.File)

			if flags.Update {
				env.Update()
			}

			table := uitable.New()
			table.MaxColWidth = 50
			table.AddRow("PRIO", "CHART", "NAME", "NAMESPACE", "RULE", "DEPLOY", "VERSION", "ACTION")

			apps, keys := env.LoadApps(envConfig, flags.Tags, flags.Apps, flags.Priorities)
			for _, key := range keys {
				for _, app := range apps[key] {
					table.AddRow(app.Declaration.Priority, app.Declaration.Helm.Chart, app.AppName, app.Namespace, app.Declaration.Helm.Version, app.CurrentVersion, app.NextVersion, app.Action.String())
				}
			}
			fmt.Println(table)
		},
	}

	addStringSliceFlag(cmd, "tag", "", []string{}, "comma separated list of tags")
	addStringSliceFlag(cmd, "priority", "", []string{}, "comma separated list of priorities")
	addStringSliceFlag(cmd, "app", "", []string{}, "application name for the action")
	addBoolFlag(cmd, "update", "", false, "update repositories before sync")

	return cmd
}
