package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func CmdEnvStart() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start environment",
		Long:  `Start existing environment.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.Name())
		},
	}

	return cmd
}
