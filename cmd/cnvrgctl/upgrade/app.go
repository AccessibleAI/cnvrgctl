package upgrade

import (
	cnvrgappv1 "github.com/cnvrgctl/pkg/cnvrgapp/api/types/v1"
	"github.com/cnvrgctl/pkg/upgrade"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var AppUpgradeCmd = &cobra.Command{
	Use:   "app",
	Short: "Execute webapp and sidekiq upgrade",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Going to run webapp and sidekiq upgrade...")
		if viper.GetBool("rollback") {
			logrus.Warnf("rolling back latest upgrade")
			upgrade.UpdateCnvrgApp(upgrade.LoadCnvrgAppFromBackup())
		} else {
			appUpgrade()
		}
		logrus.Info("done")

		//prompt := promptui.Select{
		//	Label: "Select Day",
		//	Items: []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday",
		//		"Saturday", "Sunday"},
		//}
		//
		//_, result, err := prompt.Run()
		//
		//if err != nil {
		//	fmt.Printf("Prompt failed %v\n", err)
		//	return
		//}
		//
		//fmt.Printf("You choose %q\n", result)
	},
}

func appUpgrade() {
	// check if upgrade availiable
	// check if cnvrg app deployed
	// check if cnvrg tenancy enabled
	// check if there is enough compute power for upgrade
	// get nodes
	//upgrade.GetNodes()
	if viper.GetBool("pull-app-image") {
		pullAppImage()
		//pullAppImage()
	}
	upgrade.BackupCnvrgApp()
	//upgrade.WatchForDeploymentScaleToZero()
	upgrade.SidekiqGracefulShutdown()
	//upgrade.GetNodesMetrics()
	//upgrade.GetCnvrgApp()
}

func pullAppImage() {
	cnvrgApp := upgrade.GetCnvrgApp()
	tenancyEnabled := isCnvrgTenancyEnabled(cnvrgApp)
	verifyUpgrade(cnvrgApp)
	appImage := viper.GetString("app-image")
	logrus.Infof("cnvrg tenancy enabled: %v", tenancyEnabled)
	logrus.Infof("app image for upgrade: %v", viper.GetString("app-image"))
	imagePullReady := make(chan bool)
	go upgrade.WatchForImagePullDaemonSetReady(imagePullReady)
	upgrade.DeployImagePullDaemonSet(cnvrgApp, appImage)
	<-imagePullReady
	logrus.Info("DONE")
}

func verifyUpgrade(cnvrgApp *cnvrgappv1.CnvrgApp) {
	imageEdition := getImageEdition()
	if imageEdition != cnvrgApp.Spec.CnvrgApp.Edition {
		logrus.Fatalf("deployment edition and image edition doesn't match. Deployment: %v, Image: %v ",
			cnvrgApp.Spec.CnvrgApp.Edition, imageEdition)
	}
}

func isCnvrgTenancyEnabled(app *cnvrgappv1.CnvrgApp) bool {
	if app.Spec.Tenancy.Enabled == "true" {
		return true
	}
	return false
}

func getImageEdition() string {
	if strings.Contains(viper.GetString("app-image"), "core") {
		return "core"
	}
	return "enterprise"

}
