package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mfinelli/rush/db"
	"github.com/mfinelli/rush/server"
	"github.com/mfinelli/rush/util"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Rush SSH certificate server component",
	Long: `The server component creates an HTTP server to listen (on port 8080) for
incoming certificate signing requests and then returns a signed certificate.`,
	Run: func(cmd *cobra.Command, args []string) {
		rdb, err := db.SetupDB()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		server.Serve(rdb)
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is /etc/rush/config.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("rush")

	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.auth", "htpasswd")
	viper.SetDefault("server.htpasswd", "/etc/rush/htpasswd")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("/etc/rush")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	// allow to pull nested variables from env
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	if err := util.ValidateConfig(); err != nil {
		// TODO: do better
		log.Panic(err)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
