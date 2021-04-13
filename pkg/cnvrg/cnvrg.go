package cnvrg

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	v1Core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func getK8SDefaultClient() *kubernetes.Clientset {
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
	return clientset
}

func CheckNodesReadyStatus() (bool, error) {
	ready := false
	ctx := context.Background()
	clientset := getK8SDefaultClient()
	nodes, err := clientset.CoreV1().Nodes().List(ctx, v1.ListOptions{})
	if err != nil {
		return false, err
	}
	for _, node := range nodes.Items {
		for _, condition := range node.Status.Conditions {
			if condition.Type == v1Core.NodeReady {
				ready = true
			}
		}
	}
	return ready, nil
}

//func GetCnvrgApp() (cnvrgapp *cnvrgappv1.CnvrgApp) {
//
//	config, _ := getK8SDefaultClient()
//	if err := cnvrgappv1.AddToScheme(scheme.Scheme); err != nil {
//		logrus.Debug(err.Error())
//		logrus.Fatal("Error registering cnvrgapp CR")
//	}
//	clientSet, err := cnvrgV1client.NewForConfigCnvrgApp(config)
//	if err != nil {
//		logrus.Debug(err.Error())
//		logrus.Fatal("Error creating cnvrgappv1 clientset")
//	}
//	namespace := viper.GetString("cnvrg-namespace")
//	cnvrgAppName := viper.GetString("cnvrgapp-name")
//	cnvrgapp, err = clientSet.CnvrgApps(namespace).Get(context.TODO(), cnvrgAppName, metav1.GetOptions{})
//	if err != nil {
//		logrus.Debug(err.Error())
//		logrus.Fatal("Error during fetching the cnvrgapp")
//	}
//	logrus.Debug(cnvrgapp)
//	return
//}
//
//func GetCnvrgAppUpgrade(name string) (appUpgrade *cnvrgappv1.CnvrgAppUpgrade) {
//	config, _ := getK8SDefaultClient()
//	if err := cnvrgappv1.AddToScheme(scheme.Scheme); err != nil {
//		logrus.Debug(err.Error())
//		logrus.Fatal("Error registering cnvrgapp CR")
//	}
//	clientSet, err := cnvrgV1client.NewForConfigCnvrgAppUpgrade(config)
//	if err != nil {
//		logrus.Debug(err.Error())
//		logrus.Fatal("Error creating cnvrgappv1 clientset")
//	}
//	namespace := viper.GetString("cnvrg-namespace")
//	appUpgrade, err = clientSet.CnvrgAppUpgrades(namespace).Get(context.TODO(), name, metav1.GetOptions{})
//	if err != nil {
//		logrus.Debug(err.Error())
//		logrus.Fatal("Error during fetching the cnvrgapp")
//	}
//	logrus.Debug(appUpgrade)
//	return
//}
//
//func DeleteCnvrgAppUpgrade(name string) error {
//	config, _ := getK8SDefaultClient()
//	if err := cnvrgappv1.AddToScheme(scheme.Scheme); err != nil {
//		logrus.Debug(err.Error())
//		logrus.Fatal("error registering cnvrgapp CR")
//	}
//	clientSet, err := cnvrgV1client.NewForConfigCnvrgAppUpgrade(config)
//	if err != nil {
//		logrus.Debug(err.Error())
//		return fmt.Errorf("error creating cnvrgappv1 clientset")
//	}
//	return clientSet.CnvrgAppUpgrades(viper.GetString("cnvrg-namespace")).Delete(
//		context.TODO(),
//		name,
//		metav1.DeleteOptions{})
//}
//
//func DeleteAllUpgrades() error {
//	for _, upgrade := range listAppUpgrades().Items {
//		err := DeleteCnvrgAppUpgrade(upgrade.Name)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func CreateCnvrgAppUpgrade(upgradeSpec *cnvrgappv1.CnvrgAppUpgrade) error {
//	err := ableToUpgrade()
//	if err != nil {
//		return err
//	}
//	config, _ := getK8SDefaultClient()
//	if err := cnvrgappv1.AddToScheme(scheme.Scheme); err != nil {
//		logrus.Debug(err.Error())
//		logrus.Fatal("error registering cnvrgapp CR")
//	}
//	clientSet, err := cnvrgV1client.NewForConfigCnvrgAppUpgrade(config)
//	if err != nil {
//		logrus.Debug(err.Error())
//		return fmt.Errorf("error creating cnvrgappv1 clientset")
//	}
//	res, err := clientSet.CnvrgAppUpgrades(viper.GetString("cnvrg-namespace")).Create(
//		context.TODO(),
//		upgradeSpec,
//		metav1.CreateOptions{})
//	if err != nil {
//		logrus.Debug(err.Error())
//		return fmt.Errorf("error creating upgrade spec")
//	}
//	logrus.Debug(res)
//	return nil
//}
//
//func listAppUpgrades() *cnvrgappv1.CnvrgAppUpgradeList {
//	config, _ := getK8SDefaultClient()
//	if err := cnvrgappv1.AddToScheme(scheme.Scheme); err != nil {
//		logrus.Debug(err.Error())
//		logrus.Fatal("error registering cnvrgapp CR")
//	}
//	clientSet, err := cnvrgV1client.NewForConfigCnvrgAppUpgrade(config)
//	if err != nil {
//		logrus.Debug(err.Error())
//		logrus.Fatal("error creating cnvrgappv1 clientset")
//	}
//	res, err := clientSet.CnvrgAppUpgrades(viper.GetString("cnvrg-namespace")).
//		List(context.TODO(), metav1.ListOptions{})
//	if err != nil {
//		logrus.Debug(err.Error())
//		logrus.Fatal("can't list upgrade objects")
//	}
//	logrus.Debug(res)
//	return res
//}
//
//func GetActiveAppUpgrade() *cnvrgappv1.CnvrgAppUpgrade {
//	for _, upgradeSpec := range listAppUpgrades().Items {
//		if upgradeSpec.Spec.Condition == "upgrade" {
//			return &upgradeSpec
//		}
//	}
//	return nil
//}
//
//func ableToUpgrade() error {
//	activeUpgradeName := ""
//	activeUpgradesCounts := 0
//	for _, upgradeSpec := range listAppUpgrades().Items {
//		if upgradeSpec.Spec.Condition == "upgrade" || upgradeSpec.Spec.Condition == "rollback" {
//			activeUpgradeName = upgradeSpec.Name
//			activeUpgradesCounts += 1
//		}
//	}
//	if activeUpgradesCounts > 0 {
//		return fmt.Errorf(
//			"unable create upgrade spec, upgrade: %v currently active, use --watch-upgrade to watch running upgrade", activeUpgradeName,
//		)
//	}
//	return nil
//}
//
//func WatchForCnvrgApp() {
//	//config, _ := getK8SDefaultClient()
//	//if err := cnvrgappv1.AddToScheme(scheme.Scheme); err != nil {
//	//	logrus.Debug(err.Error())
//	//	logrus.Fatal("Error registering cnvrgapp CR")
//	//}
//	//clientSet, err := cnvrgV1client.NewForConfigCnvrgApp(config)
//	//if err != nil {
//	//	logrus.Debug(err.Error())
//	//	logrus.Fatal("Error creating cnvrgappv1 clientset")
//	//}
//	//WatchCnvrgAppResources(clientSet)
//}
//
//func IsClosed(ch <-chan struct{}) bool {
//	select {
//	case <-ch:
//		return true
//	default:
//	}
//
//	return false
//}
//
//func WatchForCnvrgAppUpgrade(upgradeName string) {
//	if viper.GetBool("dry-run") {
//		logrus.Debug("Dry run on, skipping")
//		return
//	}
//	logrus.Debugf("upgradeName to watch: %v", upgradeName)
//	config, _ := getK8SDefaultClient()
//	if err := cnvrgappv1.AddToScheme(scheme.Scheme); err != nil {
//		logrus.Debug(err.Error())
//		logrus.Fatal("Error registering cnvrgapp CR")
//	}
//	clientSet, err := cnvrgV1client.NewForConfigCnvrgAppUpgrade(config)
//	if err != nil {
//		logrus.Debug(err.Error())
//		logrus.Fatal("Error creating cnvrgappv1 clientset")
//	}
//	eventData := evenHandlerData{objName: upgradeName, stopper: make(chan struct{}), messages: make(chan string, 100)}
//	s := spinner.New(spinner.CharSets[27], 50*time.Millisecond)
//	defer close(eventData.messages)
//	go pkg.StartSpinner(s, " upgrading... ", eventData.messages)
//	go WatchCnvrgAppUpgradeResources(clientSet, cnvrgAppUpgradeEventHandler, eventData)
//	<-eventData.stopper
//	s.Stop()
//}
//
//func cnvrgAppUpgradeEventHandler(eventData evenHandlerData) {
//	upgradeSpec := eventData.new.(*cnvrgappv1.CnvrgAppUpgrade)
//	if upgradeSpec.Name == eventData.objName {
//		if upgradeSpec.Spec.Condition == "upgrade" || upgradeSpec.Spec.Condition == "inactive" {
//			eventData.messages <- upgradeSpec.Name + ": " + upgradeSpec.Status.Status
//			if upgradeSpec.Status.Status == "upgrade done" {
//				close(eventData.stopper)
//			}
//		}
//	}
//}
//
//func GetAppUpgradeNameForWatch(upgradeName string) (name string) {
//	if upgradeName != "" {
//		return upgradeName
//	}
//	if viper.GetString("upgrade-name") != "" {
//		return viper.GetString("upgrade-name")
//	}
//	s := spinner.New(spinner.CharSets[27], 50*time.Millisecond)
//	go pkg.StartSpinner(s, "fetching existing upgrades...", nil)
//	var upgradeNames []string
//	for _, upgrade := range listAppUpgrades().Items {
//		upgradeNames = append(upgradeNames, upgrade.Name)
//	}
//	s.Stop()
//	prompt := promptui.Select{
//		Label: "Choose a image",
//		Items: upgradeNames,
//	}
//	_, upgradeName, err := prompt.Run()
//	if err != nil {
//		logrus.Debug(err.Error())
//		logrus.Fatal("error choosing upgrade")
//	}
//
//	return upgradeName
//}
