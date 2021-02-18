package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func CmdEnvStop() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop environment",
		Long:  `Stop existing environment.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.Name())
		},
	}

	return cmd
}
