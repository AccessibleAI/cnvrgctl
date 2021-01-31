package cnvrg

import (
	"context"
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
	objName   string
	old, new  interface{}
	stopper   chan struct{}
	messages  chan string
	eventType string
}

type updateEventHandler func(eventData evenHandlerData)

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

func WatchCnvrgAppUpgradeResources(clientSet v1.CnvrgAppUpgradeV1Interface, eventHandler updateEventHandler, eventData evenHandlerData) cache.Store {

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
				eventData.old = nil
				eventData.new = obj
				eventData.eventType = "AddFunc"
				eventHandler(eventData)
			},
			UpdateFunc: func(old, new interface{}) {
				eventData.eventType = "UpdateFunc"
				eventData.old = old
				eventData.new = new
				eventHandler(eventData)
			},
			DeleteFunc: func(obj interface{}) {
				eventData.eventType = "DeleteFunc"
				eventHandler(eventData)
			},
		},
	)

	go cnvrgAppUpgradeController.Run(eventData.stopper)
	return cnvrgAppUpgradeStore
}
