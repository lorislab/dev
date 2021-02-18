package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func CmdEnvInstall() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall environment",
		Long:  `Uninstall existing environment base from the cluster.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.Name())
		},
	}

	return cmd
}
