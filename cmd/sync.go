package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func CmdEnvSync() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync environment",
		Long:  `Sync existing environment base on the environment configuration.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.Name())
		},
	}

	return cmd
}
