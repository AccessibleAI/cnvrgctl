package main

import (
	"fmt"
	"github.com/cnvrgctl/cmd/cnvrgctl/upgrade"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "cnvrgctl",
	Short: "cnvrgctl - command line tool for managing cnvrg stack",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Setup logging
		setupLogging()
	},
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "upgrade cnvrg stack components",
}

func setupLogging() {

	// Set log verbosity
	if viper.GetBool("verbose") {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetReportCaller(true)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
		logrus.SetReportCaller(false)
	}
	// Set log format
	if viper.GetBool("json-log") {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	}
	// Logs are always goes to STDOUT
	logrus.SetOutput(os.Stdout)

}

func setupCommands() {
	// Setup commands
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "--verbose=true|false")
	rootCmd.PersistentFlags().BoolP("json-log", "J", false, "--json-log=true|false")
	upgradeCmd.PersistentFlags().BoolP("pull-app-image", "p", true, "--pull-app-image=true|false set true to pull the image on the k8s node before running the upgrade")
	upgradeCmd.AddCommand(upgrade.AppUpgradeCmd)
	rootCmd.AddCommand(upgradeCmd)
	if err := viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("json-log", rootCmd.PersistentFlags().Lookup("json-log")); err != nil {
		panic(err)
	}

}
func main() {
	setupCommands()
	//logrus.SetLevel(logrus.InfoLevel)
	x := viper.GetBool("verbose")
	if x {
		logrus.Infof("asd")
	}

	if viper.GetBool("verbose") {
		logrus.SetLevel(logrus.DebugLevel)
	}

	//logrus.Fatal("this is fatal end")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
