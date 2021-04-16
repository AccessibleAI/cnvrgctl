package main

import (
	"fmt"
	"github.com/cnvrgctl/cmd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
	"strings"
)

var rootParams = []cmd.Param{
	{Name: "verbose", Shorthand: "v", Value: false, Usage: "--verbose=true|false"},
	{Name: "json-log", Shorthand: "J", Value: false, Usage: "--json-log=true|false"},
	{Name: "cnvrgapp-name", Shorthand: "n", Value: "cnvrg-app", Usage: "cnvrgapp object name"},
	{Name: "cnvrg-namespace", Shorthand: "S", Value: "cnvrg", Usage: "cnvrgapp namespace"},
	{Name: "dry-run", Shorthand: "d", Value: false, Usage: "--dry-run=true|false"},
	{Name: "kubeconfig", Shorthand: "", Value: kubeconfigDefaultLocation(), Usage: "absolute path to the kubeconfig file"},
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

func setParams(params []cmd.Param, command *cobra.Command) {
	for _, param := range params {
		switch v := param.Value.(type) {
		case int:
			command.PersistentFlags().IntP(param.Name, param.Shorthand, v, param.Usage)
		case string:
			command.PersistentFlags().StringP(param.Name, param.Shorthand, v, param.Usage)
		case bool:
			command.PersistentFlags().BoolP(param.Name, param.Shorthand, v, param.Usage)
		}
		if err := viper.BindPFlag(param.Name, command.PersistentFlags().Lookup(param.Name)); err != nil {
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

	//setParams(cmd.ClusterUpParams, cmd.ClusterUpCmd)
	//setParams(cmd.ClusterUpParams, cmd.ClusterRemoveCmd)



	setParams(cmd.ClusterParams, cmd.ClusterCmd)
	setParams(cmd.ImagesDumpParams, cmd.DumpCmd)
	setParams(cmd.ImagesParams, cmd.ImagesCmd)

	setParams(rootParams, rootCmd)

	// cluster
	cmd.ClusterCmd.AddCommand(cmd.ClusterUpCmd)
	cmd.ClusterCmd.AddCommand(cmd.ClusterRemoveCmd)

	// images
	cmd.ImagesCmd.AddCommand(cmd.DumpCmd)
	cmd.ImagesCmd.AddCommand(cmd.PullCmd)
	cmd.ImagesCmd.AddCommand(cmd.LoadCmd)
	cmd.ImagesCmd.AddCommand(cmd.SaveCmd)
	cmd.ImagesCmd.AddCommand(cmd.TagCmd)
	cmd.ImagesCmd.AddCommand(cmd.PushCmd)

	//rootCmd.AddCommand(upgradeCmd)
	rootCmd.AddCommand(cmd.ClusterCmd)
	rootCmd.AddCommand(cmd.ImagesCmd)
	rootCmd.AddCommand(cmd.CompletionCmd)
	rootCmd.AddCommand(cmd.VersionCmd)

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
