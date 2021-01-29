package main

import (
	"github.com/cnvrgctl/pkg/cnvrgapp"
	"github.com/cnvrgctl/pkg/images"
	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var AppUpgradeCmd = &cobra.Command{
	Use:   "app",
	Short: "Execute cnvrg webapp and sidekiq upgrade",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Running cnvrg application upgrade...")
		appUpgrade()
	},
}

func appUpgrade() {
	appImage := getImageForUpgrade()
	logrus.Infof("using %v for upgrade", appImage)
}

func getImageForUpgrade() string {
	appImage := viper.GetString("app-image")
	if appImage != "" {
		return appImage
	}
	cnvrgSpec := cnvrgapp.GetCnvrgApp()
	logrus.Info(cnvrgSpec)
	prompt := promptui.Select{
		Label: "Choose a image",
		Items: images.ListAppImages(),
	}
	_, appImage, err := prompt.Run()
	if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("error choosing image for upgrade")
	}
	return appImage
}

