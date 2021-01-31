package main

import (
	"encoding/json"
	"github.com/cnvrgctl/pkg/cnvrg"
	v1 "github.com/cnvrgctl/pkg/cnvrg/api/types/v1"
	"github.com/cnvrgctl/pkg"
	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var upgradeAppParams = []param{
	{name: "condition", value: "upgrade", usage: "upgrade | rollback | inactive"},
	{name: "cacheDsName", value: "app-image-cache", usage: "caching DaemonSet name"},
	{name: "cnvrgAppName", value: "cnvrg-app", usage: "cnvrgapp object name"},
	{name: "image", value: "", usage: "image for upgrade"},
	{name: "cacheImage", value: "true", usage: "true/false to cache image before upgrade"},
	{name: "watch-upgrade", value: false, usage: "--watch-upgrade=true to watch for existing upgrade"},
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "upgrade cnvrg stack components",
}

var appUpgradeCmd = &cobra.Command{
	Use:   "app",
	Short: "Execute cnvrg webapp and sidekiq upgrade",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("running cnvrg application upgrade...")
		appUpgrade()
	},
}

func appUpgrade() {

	if viper.GetBool("watch-upgrade") == false {
		appImage := getImageForUpgrade()
		logrus.Infof("image: %v", appImage)
		upgradeSpec := v1.NewCnvrgAppUpgrade(appImage)
		if viper.GetBool("dry-run") {
			b, _ := json.MarshalIndent(upgradeSpec, "", "  ")
			logrus.Info("\n" + string(b))
		}
		cnvrg.CreateCnvrgAppUpgrade(upgradeSpec)
	}
	cnvrg.WatchForCnvrgAppUpgrade()
}

func getImageForUpgrade() string {
	appImage := viper.GetString("app-image")
	if appImage != "" {
		return appImage
	}
	cnvrgSpec := cnvrg.GetCnvrgApp()
	logrus.Debug(cnvrgSpec)
	if cnvrgSpec.Spec.CnvrgApp.Conf.Registry.URL != "docker.io" {
		logrus.Fatalf("can't list images, docker registry: %v not supported. explicitly provide app image with --app-image flag",
			cnvrgSpec.Spec.CnvrgApp.Conf.Registry.URL)
	}
	if cnvrgSpec.Spec.CnvrgApp.Conf.Registry.User == "" || cnvrgSpec.Spec.CnvrgApp.Conf.Registry.Password == "" {
		logrus.Fatal("can't list images, missing registry credentials. explicitly provide app image with --app-image flag",
			cnvrgSpec.Spec.CnvrgApp.Conf.Registry.URL)
	}
	prompt := promptui.Select{
		Label: "Choose a image",
		Items: pkg.ListAppImages(
			cnvrgSpec.Spec.CnvrgApp.Conf.Registry.User,
			cnvrgSpec.Spec.CnvrgApp.Conf.Registry.Password,
		),
	}
	_, appImage, err := prompt.Run()
	if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("error choosing image for upgrade")
	}
	return appImage
}
