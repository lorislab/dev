package cmd

import (
	"io/ioutil"

	"github.com/lorislab/dev/env"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

//Env environment commands
func Env() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "env",
		Short: "Environment commands",
		Long:  `Environment commands for the envorinment.`,
	}

	addFlag(cmd, "env-config", "e", "env.yaml", "Environment configuration")

	addChildCmd(cmd, EnvSync())
	addChildCmd(cmd, EnvStatus())
	addChildCmd(cmd, EnvUninstall())
	return cmd
}

type envFlags struct {
	File       string   `mapstructure:"env-config"`
	Apps       []string `mapstructure:"app"`
	Tags       []string `mapstructure:"tag"`
	Priorities []string `mapstructure:"priority"`
	Update     bool     `mapstructure:"update"`
}

func envConfig(file string) *env.LocalEnvironment {
	clusterConfig := &env.LocalEnvironment{}
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal().Str("file", file).Err(err).Msg("Error loading the file")
	}
	err = yaml.Unmarshal(yamlFile, &clusterConfig)
	if err != nil {
		panic(err)
	}
	return clusterConfig
}

func readEnvFlags() envFlags {
	options := envFlags{}
	readOptions(&options)
	return options
}
