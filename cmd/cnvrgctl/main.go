package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "cnvrgctl",
	Short: "cnvrgctl - command line tool for managing cnvrg stack",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// app-image is required only if rollback flag not set
		//if !viper.GetBool("rollback") {
		//	if err := upgradeCmd.MarkFlagRequired("app-image"); err != nil {
		//		panic(err)
		//	}
		//}
		// Setup logging
		setupLogging()
		logrus.Debugf("kubeconfig: %v", viper.GetString("kubeconfig"))
		logrus.Debugf("verbose: %v", viper.GetBool("verbose"))
		logrus.Debugf("json-log: %v", viper.GetBool("json-log"))
		logrus.Debugf("cache-image: %v", viper.GetBool("cache-image"))
		logrus.Debugf("cnvrgapp-name: %v", viper.GetBool("cnvrgapp-name"))
		logrus.Debugf("cnvrgapp-name: %v", viper.GetBool("cnvrgapp-name"))
		logrus.Debugf("rollback: %v", viper.GetBool("rollback"))
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
	// Init config
	cobra.OnInitialize(initConfig)
	// Setup commands
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "--verbose=true|false")
	rootCmd.PersistentFlags().BoolP("json-log", "J", false, "--json-log=true|false")
	rootCmd.PersistentFlags().StringP("cnvrgapp-name", "n", "cnvrg-app", "name of the CnvrgApp spec")
	rootCmd.PersistentFlags().StringP("cnvrg-namespace", "S", "cnvrg", "CnvrgApp namespace")
	rootCmd.PersistentFlags().BoolP("dry-run", "d", false, "--dry-run=true|false")
	upgradeCmd.PersistentFlags().BoolP("cache-image", "c", true, "--cache-image=true|false set true to pull the image on the k8s node before running the upgrade")
	upgradeCmd.PersistentFlags().StringP("app-image", "i", "", "app image to use for upgrade")
	upgradeCmd.PersistentFlags().BoolP("rollback", "r", false, "rollback to previous cnvrgapp")

	kubeconfigDefaultLocation := ""
	if home := homedir.HomeDir(); home != "" {
		kubeconfigDefaultLocation = filepath.Join(home, ".kube", "config")
	}
	rootCmd.PersistentFlags().String("kubeconfig", kubeconfigDefaultLocation, "absolute path to the kubeconfig file")
	upgradeCmd.AddCommand(AppUpgradeCmd)
	rootCmd.AddCommand(upgradeCmd)
	if err := viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("json-log", rootCmd.PersistentFlags().Lookup("json-log")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("kubeconfig", rootCmd.PersistentFlags().Lookup("kubeconfig")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("cnvrgapp-name", rootCmd.PersistentFlags().Lookup("cnvrgapp-name")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("cnvrg-namespace", rootCmd.PersistentFlags().Lookup("cnvrg-namespace")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("dry-run", rootCmd.PersistentFlags().Lookup("dry-run")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("cache-image", upgradeCmd.PersistentFlags().Lookup("cache-image")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("pull-app-image", upgradeCmd.PersistentFlags().Lookup("pull-app-image")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("app-image", upgradeCmd.PersistentFlags().Lookup("app-image")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("rollback", upgradeCmd.PersistentFlags().Lookup("rollback")); err != nil {
		panic(err)
	}

}

func initConfig() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
}

func main() {
	setupCommands()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
