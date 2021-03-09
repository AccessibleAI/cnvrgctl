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

type param struct {
	name      string
	shorthand string
	value     interface{}
	usage     string
	required  bool
}

var rootParams = []param{
	{name: "verbose", shorthand: "v", value: false, usage: "--verbose=true|false"},
	{name: "json-log", shorthand: "J", value: false, usage: "--json-log=true|false"},
	{name: "cnvrgapp-name", shorthand: "n", value: "cnvrg-app", usage: "cnvrgapp object name"},
	{name: "cnvrg-namespace", shorthand: "S", value: "cnvrg", usage: "cnvrgapp namespace"},
	{name: "dry-run", shorthand: "d", value: false, usage: "--dry-run=true|false"},
	{name: "kubeconfig", shorthand: "", value: kubeconfigDefaultLocation(), usage: "absolute path to the kubeconfig file"},
}

var rootCmd = &cobra.Command{
	Use:   "cnvrgctl",
	Short: "cnvrgctl - command line tool for managing cnvrg stack",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		setupLogging()
	},
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

func setParams(params []param, command *cobra.Command) {
	for _, param := range params {
		switch v := param.value.(type) {
		case int:
			command.PersistentFlags().IntP(param.name, param.shorthand, v, param.usage)
		case string:
			command.PersistentFlags().StringP(param.name, param.shorthand, v, param.usage)
		case bool:
			command.PersistentFlags().BoolP(param.name, param.shorthand, v, param.usage)
		}
		if err := viper.BindPFlag(param.name, command.PersistentFlags().Lookup(param.name)); err != nil {
			panic(err)
		}
	}
}

func setupCommands() {
	// Init config
	cobra.OnInitialize(initConfig)
	// update cmd
	//setParams(upgradeAppParams, appUpgradeCmd)
	//upgradeCmd.AddCommand(appUpgradeCmd)

	setParams(imagesDumpParams, dumpCmd)
	setParams(imagesParams, imagesCmd)
	setParams(rootParams, rootCmd)

	imagesCmd.AddCommand(dumpCmd)
	imagesCmd.AddCommand(pullCmd)
	imagesCmd.AddCommand(loadCmd)
	imagesCmd.AddCommand(saveCmd)
	imagesCmd.AddCommand(tagCmd)
	imagesCmd.AddCommand(pushCmd)

	//rootCmd.AddCommand(upgradeCmd)
	rootCmd.AddCommand(imagesCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(versionCmd)

}

func kubeconfigDefaultLocation() string {
	kubeconfigDefaultLocation := ""
	if home := homedir.HomeDir(); home != "" {
		kubeconfigDefaultLocation = filepath.Join(home, ".kube", "config")
	}
	return kubeconfigDefaultLocation
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
