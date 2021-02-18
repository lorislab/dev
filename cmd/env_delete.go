package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func CmdEnvDelete() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete environment",
		Long:  `Delete exists environment.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.Name())
		},
	}

	return cmd
}
