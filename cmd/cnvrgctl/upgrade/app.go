package upgrade

import (
	"github.com/cnvrgctl/pkg/k8s"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var AppUpgradeCmd = &cobra.Command{
	Use:   "app",
	Short: "Execute webapp and sidekiq upgrade",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("going to run webapp and sidekiq upgrade")

		appUpgrade()

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
	// check if k8s availiable
	// check if cnvrg app deployed
	// check if cnvrg tenancy enabled
	// check if there is enough compute power for upgrade
	// get nodes
	//k8s.GetNodes()
	k8s.GetCnvrgApp()
}
