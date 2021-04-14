package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var CnvrgUpParams = []Param{

}

var CnvrgCmd = &cobra.Command{
	Use:   "cnvrg",
	Short: "deploy cnvrg stack",
}

var CnvrgUpCmd = &cobra.Command{
	Use:   "up",
	Short: "bring up cnvrg stack",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("installing cnvrg stack")
	},
}


