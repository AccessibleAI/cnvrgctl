package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Version string = "v1.0.1"



var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print cnvrgctl version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("version : " + Version)
	},
}
