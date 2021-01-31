package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	buildVersion string = ""
	commit       string = ""
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print cnvrgctl version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("version : " + buildVersion)
		fmt.Println("commit  : " + commit)
	},
}
