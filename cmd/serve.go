package cmd

import (
	"github.com/mfinelli/rush/server"
)

import "github.com/spf13/cobra"

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Rush SSH certificate server component",
	Long: `The server component creates an HTTP server to listen (on port 8080) for
incoming certificate signing requests and then returns a signed certificate.`,
	Run: func(cmd *cobra.Command, args []string) {
		server.Serve()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
