package cnvrg

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/cnvrgctl/pkg"
	cnvrgappv1 "github.com/cnvrgctl/pkg/cnvrg/api/types/v1"
	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	cnvrgV1client "github.com/cnvrgctl/pkg/cnvrg/clientset/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"context"
	"time"
)

func getK8SDefaultClient() (*rest.Config, *kubernetes.Clientset) {
	config, err := clientcmd.BuildConfigFromFlags("", viper.GetString("kubeconfig"))
	if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("Error building kubeconfig")
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("Error creating client")
	}
	return config, clientset
}

func GetCnvrgApp() (cnvrgapp *cnvrgappv1.CnvrgApp) {

	config, _ := getK8SDefaultClient()
	if err := cnvrgappv1.AddToScheme(scheme.Scheme); err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("Error registering cnvrgapp CR")
	}
	clientSet, err := cnvrgV1client.NewForConfigCnvrgApp(config)
	if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("Error creating cnvrgappv1 clientset")
	}
	namespace := viper.GetString("cnvrg-namespace")
	cnvrgAppName := viper.GetString("cnvrgapp-name")
	cnvrgapp, err = clientSet.CnvrgApps(namespace).Get(context.TODO(), cnvrgAppName, metav1.GetOptions{})
	if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("Error during fetching the cnvrgapp")
	}
	logrus.Debug(cnvrgapp)
	return
}

func CreateCnvrgAppUpgrade(upgradeSpec *cnvrgappv1.CnvrgAppUpgrade) {
	ok, msg := ableToUpgrade()
	if ok == false {
		logrus.Fatal(msg)
	}

	config, _ := getK8SDefaultClient()
	if err := cnvrgappv1.AddToScheme(scheme.Scheme); err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("error registering cnvrgapp CR")
	}
	clientSet, err := cnvrgV1client.NewForConfigCnvrgAppUpgrade(config)
	if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("error creating cnvrgappv1 clientset")
	}
	res, err := clientSet.CnvrgAppUpgrades(viper.GetString("cnvrg-namespace")).Create(
		context.TODO(),
		upgradeSpec,
		metav1.CreateOptions{})
	if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("error creating upgrade spec")
	}
	logrus.Debug(res)
}

func listAppUpgrades() *cnvrgappv1.CnvrgAppUpgradeList {
	config, _ := getK8SDefaultClient()
	if err := cnvrgappv1.AddToScheme(scheme.Scheme); err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("error registering cnvrgapp CR")
	}
	clientSet, err := cnvrgV1client.NewForConfigCnvrgAppUpgrade(config)
	if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("error creating cnvrgappv1 clientset")
	}
	res, err := clientSet.CnvrgAppUpgrades(viper.GetString("cnvrg-namespace")).
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("can't list upgrade objects")
	}
	logrus.Debug(res)
	return res
}

func GetActiveAppUpgrade() *cnvrgappv1.CnvrgAppUpgrade {
	for _, upgradeSpec := range listAppUpgrades().Items {
		if upgradeSpec.Spec.Condition == "upgrade" {
			return &upgradeSpec
		}
	}
	return nil
}

func ableToUpgrade() (bool, string) {
	activeUpgradeName := ""
	activeUpgradesCounts := 0
	for _, upgradeSpec := range listAppUpgrades().Items {
		if upgradeSpec.Spec.Condition == "upgrade" || upgradeSpec.Spec.Condition == "rollback" {
			activeUpgradeName = upgradeSpec.Name
			activeUpgradesCounts += 1
		}
	}
	if activeUpgradesCounts > 0 {
		return false, fmt.Sprintf(
			"unable create upgrade spec, upgrade: %v currently active, use --watch-upgrade to watch running upgrade", activeUpgradeName,
		)
	}
	return true, ""
}

func WatchForCnvrgApp() {
	//config, _ := getK8SDefaultClient()
	//if err := cnvrgappv1.AddToScheme(scheme.Scheme); err != nil {
	//	logrus.Debug(err.Error())
	//	logrus.Fatal("Error registering cnvrgapp CR")
	//}
	//clientSet, err := cnvrgV1client.NewForConfigCnvrgApp(config)
	//if err != nil {
	//	logrus.Debug(err.Error())
	//	logrus.Fatal("Error creating cnvrgappv1 clientset")
	//}
	//WatchCnvrgAppResources(clientSet)
}

func IsClosed(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

func WatchForCnvrgAppUpgrade(upgradeName string) {
	if viper.GetBool("dry-run") {
		logrus.Debug("Dry run on, skipping")
		return
	}
	logrus.Debugf("upgradeName to watch: %v", upgradeName)
	config, _ := getK8SDefaultClient()
	if err := cnvrgappv1.AddToScheme(scheme.Scheme); err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("Error registering cnvrgapp CR")
	}
	clientSet, err := cnvrgV1client.NewForConfigCnvrgAppUpgrade(config)
	if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("Error creating cnvrgappv1 clientset")
	}
	eventData := evenHandlerData{objName: upgradeName, stopper: make(chan struct{}), messages: make(chan string, 100)}
	s := spinner.New(spinner.CharSets[27], 50*time.Millisecond)
	defer close(eventData.messages)
	go pkg.StartSpinner(s, " upgrading... ", eventData.messages)
	go WatchCnvrgAppUpgradeResources(clientSet, cnvrgAppUpgradeEventHandler, eventData)
	<-eventData.stopper
	s.Stop()
}

func cnvrgAppUpgradeEventHandler(eventData evenHandlerData) {
	upgradeSpec := eventData.new.(*cnvrgappv1.CnvrgAppUpgrade)
	if upgradeSpec.Name == eventData.objName {
		if upgradeSpec.Spec.Condition == "upgrade" || upgradeSpec.Spec.Condition == "inactive" {
			eventData.messages <- upgradeSpec.Name + ": " + upgradeSpec.Status.Status
			if upgradeSpec.Status.Status == "upgrade done" {
				close(eventData.stopper)
			}
		}
	}
}

func GetAppUpgradeNameForWatch(upgradeName string) (name string) {
	if upgradeName != "" {
		return upgradeName
	}
	if viper.GetString("upgrade-name") != "" {
		return viper.GetString("upgrade-name")
	}
	s := spinner.New(spinner.CharSets[27], 50*time.Millisecond)
	go pkg.StartSpinner(s, "fetching existing upgrades...", nil)
	var upgradeNames []string
	for _, upgrade := range listAppUpgrades().Items {
		upgradeNames = append(upgradeNames, upgrade.Name)
	}
	s.Stop()
	prompt := promptui.Select{
		Label: "Choose a image",
		Items: upgradeNames,
	}
	_, upgradeName, err := prompt.Run()
	if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("error choosing upgrade")
	}

	return upgradeName
}
