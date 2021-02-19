package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	goVersion "go.hein.dev/go-version"
)

var (
	name = "dev"
	// Used for flags.
	shortened = false
	output    = "json"
	bv        BuildVersion
	cfgFile   string
	v         string
	rootCmd   = &cobra.Command{
		Use:   name,
		Short: "Developer tool for k8s development",
		Long:  `Developer tool for k8s with helm chart development.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := setUpLogs(os.Stdout, v); err != nil {
				return err
			}
			return nil
		},
	}
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Version will output the current build information",
		Long:  ``,
		Run: func(_ *cobra.Command, _ []string) {
			resp := goVersion.FuncWithOutput(shortened, bv.Version, bv.Commit, bv.Date, output)
			fmt.Print(resp)
		},
	}
)

//BuildVersion build version
type BuildVersion struct {
	Version string
	Commit  string
	Date    string
}

// Execute executes the root command.
func Execute(version BuildVersion) {
	bv = version
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal().Err(err).Msg("Error execute command")
	}
}

func init() {
	versionCmd.Flags().BoolVarP(&shortened, "short", "s", false, "Print just the version number.")
	versionCmd.Flags().StringVarP(&output, "output", "o", "json", "Output format. One of 'yaml' or 'json'.")
	rootCmd.AddCommand(versionCmd)

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/."+name+".yaml)")
	rootCmd.PersistentFlags().StringVarP(&v, "level", "v", zerolog.InfoLevel.String(), "Log level (debug, info, warn, error, fatal, panic)")

	rootCmd.AddCommand(Env())
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal().Err(err).Msg("Error find homer directory for the user")
		}

		// Search config in home directory with name ".dev" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName("." + name)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetEnvPrefix(strings.ToUpper(name))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Debug().Str("file", viper.ConfigFileUsed()).Msg("Using config")
	}
}

func setUpLogs(out io.Writer, level string) error {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(lvl)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: out})
	return nil
}
