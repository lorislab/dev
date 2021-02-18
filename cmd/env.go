package cmd

import (
	"io/ioutil"

	"github.com/lorislab/dev/dev"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func CmdEnv() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "env",
		Short: "Environment commands",
		Long:  `Environment commands for the envorinment.`,
	}

	addFlag(cmd, "env-config", "e", "env.yaml", "Environment configuration")

	addChildCmd(cmd, CmdEnvStart())
	addChildCmd(cmd, CmdEnvStop())
	addChildCmd(cmd, CmdEnvCreate())
	addChildCmd(cmd, CmdEnvDelete())
	addChildCmd(cmd, CmdEnvSync())
	addChildCmd(cmd, CmdEnvStatus())
	return cmd
}

type envFlags struct {
	File string `mapstructure:"env-config"`
}

type appFlags struct {
	Env        envFlags `mapstructure:",squash"`
	Apps       []string `mapstructure:"app"`
	Tags       []string `mapstructure:"tag"`
	Priorities []string `mapstructure:"priority"`
	Update     bool     `mapstructure:"update"`
}

var appActionStr = []string{
	"",
	"",
	"install",
	"upgrade",
	"downgrade",
	"uninstall",
}

func envConfig(env envFlags) dev.LocalEnvironment {
	clusterConfig := dev.LocalEnvironment{}
	yamlFile, err := ioutil.ReadFile(env.File)
	if err != nil {
		log.WithFields(log.Fields{
			"file":  env.File,
			"error": err,
		}).Fatal("Error loading the file")
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

func readAppFlags() appFlags {
	options := appFlags{}
	readOptions(&options)
	return options
}
