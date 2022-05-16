package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"

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
	rootCmd.AddCommand(serveCmd)
}
