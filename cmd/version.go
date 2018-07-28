package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionString = "undefined"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of signmykey",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(versionString)
	},
}
