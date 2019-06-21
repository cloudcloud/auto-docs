// Package main is where the base binary will be built from.
//
// From this base location, sub-directories are used to capture
// models and business logic as required, grouped into their
// appropriate locations.
//
// The binary that is built from this package will begin an
// auto-docs process that encompasses a HTTP server, a data
// fetcher, and an internal data store.
package main

import (
	"log"
	"os"

	autodocs "github.com/cloudcloud/auto-docs"
	"github.com/cloudcloud/auto-docs/auto-docs/server"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// GitCommit allows for the assignment of the current git
	// commit to use when displaying version information.
	GitCommit = ""

	// cfgFile is a reference to our current configuration file.
	cfgFile = ""

	// rootCmd is the base command our tree is built from.
	rootCmd = &cobra.Command{}

	// config has an instance of our local config setup.
	config = &autodocs.Config{}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd = &cobra.Command{Short: "Gateway to the world of auto docs.", ValidArgs: []string{"version", "help", "server"}}
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")

	viper.SetDefault("Git.SSHKey", "/var/auto-docs/keys/id_rsa")
	viper.SetDefault("Git.Branch", "master")
	viper.SetDefault("Git.LocalPath", "/tmp/auto-docs")
	viper.SetDefault("Git.Timeout", "1500")
	viper.SetDefault("Git.Period", "300")
	viper.SetDefault("Listen", ":9003")
	viper.SetDefault("Name", "auto-docs")
}

func main() {
	a := os.Args
	rootCmd.Use = a[0]

	cobra.OnlyValidArgs(rootCmd, a[1:])
	rootCmd.AddCommand(buildVersionCommand())
	rootCmd.AddCommand(buildServerCommand())

	rootCmd.SetArgs(a[1:])
	rootCmd.Execute()
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		h, err := homedir.Dir()
		if err != nil {
			log.Fatalf("%s", err)
		}

		viper.AddConfigPath(h)
		viper.SetConfigName(".ad")
	}

	if err := viper.ReadInConfig(); err != nil {
		viper.WriteConfig()
	}

	viper.Unmarshal(&config)
}

func buildVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Current version of auto-docs",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("auto-docs", GitCommit)
		},
	}
}

func buildServerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Start the auto-docs server",
		Run: func(cmd *cobra.Command, args []string) {
			server.New(config).Start()
		},
	}
}
