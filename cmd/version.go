package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mfinelli/rush/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the program version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("rush version %s\n", version.VERSION.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
