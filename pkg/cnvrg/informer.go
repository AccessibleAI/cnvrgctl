package cnvrg

import (
	"context"
	"github.com/briandowns/spinner"
	cnvrgappv1 "github.com/cnvrgctl/pkg/cnvrg/api/types/v1"
	v1 "github.com/cnvrgctl/pkg/cnvrg/clientset/v1"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"time"
)

type evenHandlerData struct {
	old, new interface{}
	stopper  chan struct{}
	spinner  spinner.Spinner
}

type updateEventHandler func(old, new interface{}, stopper chan struct{}, message chan string)

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
			},
			UpdateFunc: func(old, new interface{}) {
			},
			DeleteFunc: func(obj interface{}) {
			},
		},
	)

	go cnvrgAppController.Run(wait.NeverStop)
	return cnvrgAppStore
}

func WatchCnvrgAppUpgradeResources(clientSet v1.CnvrgAppUpgradeV1Interface, eventHandler updateEventHandler, stopper chan struct{}, message chan string) cache.Store {

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
				eventHandler(nil, obj, stopper, message)
			},
			UpdateFunc: func(old, new interface{}) {
				eventHandler(old, new, stopper, message)
			},
			DeleteFunc: func(obj interface{}) {
				eventHandler(nil, obj, stopper, message)
			},
		},
	)
	go cnvrgAppUpgradeController.Run(stopper)
	return cnvrgAppUpgradeStore
}
