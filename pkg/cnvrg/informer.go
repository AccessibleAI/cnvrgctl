package cnvrg

import (
	"context"
	cnvrgappv1 "github.com/cnvrgctl/pkg/cnvrg/api/types/v1"
	v1 "github.com/cnvrgctl/pkg/cnvrg/clientset/v1"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"time"
)

func WatchCnvrgAppResources(clientSet v1.CnvrgAppV1Interface) cache.Store {

	cnvrgAppStore, cnvrgAppController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return clientSet.CnvrgApps(viper.GetString("cnvrg-namespace")).List(context.TODO(), lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return clientSet.CnvrgApps(viper.GetString("cnvrg-namespace")).Watch(context.TODO(), lo)
			},
		},
		&cnvrgappv1.CnvrgApp{},
		10*time.Second,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				logrus.Info("cnvrgapp has ben created")
			},
			UpdateFunc: func(old, new interface{}) {
				logrus.Info("The cnvrgapp has been updated")
			},
			DeleteFunc: func(obj interface{}) {
				logrus.Info("the cnvrgapp has been deleted")
			},
		},
	)

	go cnvrgAppController.Run(wait.NeverStop)
	return cnvrgAppStore
}

func WatchCnvrgAppUpgradeResources(clientSet v1.CnvrgAppUpgradeV1Interface) cache.Store {

	cnvrgAppUpgradeStore, cnvrgAppUpgradeController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return clientSet.CnvrgAppUpgrades(viper.GetString("cnvrg-namespace")).List(context.TODO(), lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return clientSet.CnvrgAppUpgrades(viper.GetString("cnvrg-namespace")).Watch(context.TODO(), lo)
			},
		},
		&cnvrgappv1.CnvrgAppUpgrade{},
		10*time.Second,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				logrus.Info("upgrade has ben created")
			},
			UpdateFunc: func(old, new interface{}) {
				logrus.Info("the upgrade has been updated")
			},
			DeleteFunc: func(obj interface{}) {
				logrus.Info("the upgrade has been deleted")
			},
		},
	)

	go cnvrgAppUpgradeController.Run(wait.NeverStop)
	return cnvrgAppUpgradeStore
}
