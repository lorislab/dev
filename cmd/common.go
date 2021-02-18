package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func addChildCmd(parent, child *cobra.Command) {
	parent.AddCommand(child)
	child.Flags().AddFlagSet(parent.Flags())
}

func addBoolFlag(command *cobra.Command, name, shorthand string, value bool, usage string) *pflag.Flag {
	command.Flags().BoolP(name, shorthand, value, usage)
	return addViper(command, name)
}

func addStringSliceFlag(command *cobra.Command, name, shorthand string, value []string, usage string) *pflag.Flag {
	command.Flags().StringSliceP(name, shorthand, value, usage)
	return addViper(command, name)
}

func addFlag(command *cobra.Command, name, shorthand string, value string, usage string) *pflag.Flag {
	command.Flags().StringP(name, shorthand, value, usage)
	return addViper(command, name)
}

func addViper(command *cobra.Command, name string) *pflag.Flag {
	f := command.Flags().Lookup(name)
	err := viper.BindPFlag(name, f)
	if err != nil {
		log.WithField("name", name).Panic(err)
	}
	return f
}

func readOptions(options interface{}) {
	err := viper.Unmarshal(options)
	if err != nil {
		log.Panic(err)
	}
	log.WithField("options", options).Debug("Load options")
}
