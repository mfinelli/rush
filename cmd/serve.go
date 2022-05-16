package cmd

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mfinelli/rush/db"
	"github.com/mfinelli/rush/server"
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

	// env variables get parsed as strings, but we use some as ints
	// check those types and convert to int if necessary
	if reflect.TypeOf(viper.Get("server.port")).Kind() == reflect.String {
		intval, err := strconv.Atoi(viper.Get("server.port").(string))
		if err != nil {
			// TODO: we should handle this better
			//       if we don't get a valid port (<65k) then we
			//       need to let the user know
			log.Panic(err)
		}

		viper.Set("server.port", intval)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
