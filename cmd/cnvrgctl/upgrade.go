package main

import (
	"github.com/cnvrgctl/pkg/cnvrgapp"
	v1 "github.com/cnvrgctl/pkg/cnvrgapp/api/types/v1"
	"github.com/cnvrgctl/pkg/images"
	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"encoding/json"
)

var AppUpgradeCmd = &cobra.Command{
	Use:   "app",
	Short: "Execute cnvrg webapp and sidekiq upgrade",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("running cnvrg application upgrade...")
		appUpgrade()
	},
}

func appUpgrade() {
	appImage := getImageForUpgrade()
	logrus.Infof("using %v for upgrade", appImage)
	upgradeSpec := v1.NewCnvrgAppUpgrade(
		viper.GetString("cnvrg-namespace"),
		viper.GetString("cnvrgapp-name"),
		appImage,
		viper.GetString("cache-image"),
	)
	if viper.GetBool("dry-run") {
		b, _ := json.MarshalIndent(upgradeSpec, "", "  ")
		logrus.Info("\n" + string(b))
	}
	cnvrgapp.CreateCnvrgAppUpgrade(upgradeSpec)

}

func getImageForUpgrade() string {
	appImage := viper.GetString("app-image")
	if appImage != "" {
		return appImage
	}
	cnvrgSpec := cnvrgapp.GetCnvrgApp()
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
		Items: images.ListAppImages(
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
