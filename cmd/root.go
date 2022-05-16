package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rush",
	Short: "An SSH certificate client and server",
	Long: `Rush is an SSH certificate client and server. The server component will
issue a signed, short-lived user SSH certificate, or a server host
certificate. The client can be used to authenticate and ask the server
for a either a user or host certificate.`,
}

// Execute adds all child commands to the root command and sets flags
// appropriately. This is called by main.main(). It only needs to happen once
// to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
