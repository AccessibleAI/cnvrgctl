package cnvrg

import (
	"fmt"
	cnvrgappv1 "github.com/cnvrgctl/pkg/cnvrg/api/types/v1"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	cnvrgV1client "github.com/cnvrgctl/pkg/cnvrg/clientset/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"context"
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

func ableToUpgrade() (bool, string) {
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
	activeUpgradeName := ""
	activeUpgradesCounts := 0
	for _, upgradeSpec := range res.Items {
		if upgradeSpec.Spec.Condition == "upgrade" || upgradeSpec.Spec.Condition == "rollback" {
			activeUpgradeName = upgradeSpec.Name
			activeUpgradesCounts += 1
		}
	}

	if activeUpgradesCounts > 0 {
		return false, fmt.Sprintf("unable create upgrade spec, upgrade: %v currently active", activeUpgradeName)
	}

	return true, ""
}
