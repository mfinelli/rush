package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var VERSION string = "1.0.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the program version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("rush version %s\n", VERSION)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
