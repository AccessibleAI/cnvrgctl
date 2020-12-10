package k8s

import (
	"context"
	"fmt"
	cnvrgappinformer "github.com/cnvrgctl/pkg/cnvrgapp"
	cnvrgappv1 "github.com/cnvrgctl/pkg/cnvrgapp/api/types/v1"
	cnvrgappV1client "github.com/cnvrgctl/pkg/cnvrgapp/clientset/v1"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
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

func GetNodes() *v1.NodeList {
	_, client := getK8SDefaultClient()
	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Errorf("Can't fetch K8S cluster nodes")
		logrus.Debug(err.Error())
	}
	return nodes
}

func GetCnvrgApp() (cnvrgapp *cnvrgappv1.CnvrgApp) {

	config, _ := getK8SDefaultClient()
	if err := cnvrgappv1.AddToScheme(scheme.Scheme); err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("Error registering cnvrgapp CR")
	}
	clientSet, err := cnvrgappV1client.NewForConfig(config)
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
	store := cnvrgappinformer.WatchResources(clientSet)


	appObj, exists, err := store.GetByKey("cnvrg/cnvrg-app")

	if exists {
		app := appObj.(*cnvrgappv1.CnvrgApp)
		logrus.Info(app.Name)
	}

	return

}

func GetAppDeployment(webAppDeploymentName string) (appDeploy *apps.Deployment) {
	logrus.Info("Getting webapp deployment")
	_, client := getK8SDefaultClient()
	clientset := client.AppsV1().Deployments(viper.GetString("cnvrg-namespace"))
	appDeploy, err := clientset.Get(context.TODO(), webAppDeploymentName, metav1.GetOptions{})
	if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("Error fetching app deployment")
	}
	return

}

func DeployImagePullDaemonSet(cnvrgApp *cnvrgappv1.CnvrgApp, image string) {
	logrus.Info("starting image pull daemon set...")
	dsName := "cnvrg-image-puller"
	labelSet := map[string]string{"ds-name": dsName}
	command := []string{"/bin/bash", "-c", "sleep inf"}
	appDeployment := GetAppDeployment(cnvrgApp.Spec.CnvrgApp.SvcName)
	logrus.Debugf("image cache ds using pull secret: %v", appDeployment.Spec.Template.Spec.ImagePullSecrets)
	logrus.Debugf("image cache ds using toleration: %v", appDeployment.Spec.Template.Spec.Tolerations)
	logrus.Debugf("image cache ds using node selector: %v", appDeployment.Spec.Template.Spec.NodeSelector)
	ds := &apps.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: viper.GetString("cnvrg-namespace"),
			Name:      dsName,
		},
		Spec: apps.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labelSet,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labelSet,
				},
				Spec: v1.PodSpec{
					ImagePullSecrets: appDeployment.Spec.Template.Spec.ImagePullSecrets,
					Tolerations:      appDeployment.Spec.Template.Spec.Tolerations,
					NodeSelector:     appDeployment.Spec.Template.Spec.NodeSelector,
					Containers: []v1.Container{
						{
							Name:    dsName,
							Image:   image,
							Command: command,
						},
					},
				},
			},
		},
	}
	logrus.Debugf("Gonna create pull image DaemonSet: %v", ds)
	_, client := getK8SDefaultClient()
	dsClientSet := client.AppsV1().DaemonSets(viper.GetString("cnvrg-namespace"))
	_, err := dsClientSet.Create(context.TODO(), ds, metav1.CreateOptions{})
	if errors.IsAlreadyExists(err) {
		logrus.Warn("the Image Pull DaemonSet already exists, gonna update it")
		_, err = dsClientSet.Update(context.TODO(), ds, metav1.UpdateOptions{})
		if err != nil {
			logrus.Debug(err.Error())
			logrus.Fatal("Wasn't able to update Image Pull DaemonSet")
		}
	} else if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("Wasn't able to create Image Pull DaemonSet ")
	}
	logrus.Info("The Image Pull DaemonSet successfully created")

	factory := informers.NewSharedInformerFactory(client, 0)
	informer := factory.Apps().V1().DaemonSets().Informer()
	stopper := make(chan struct{})
	defer close(stopper)
	defer runtime.HandleCrash()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onAdd,
		UpdateFunc: onUpdate,
	})
	go informer.Run(stopper)
	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}
	<-stopper
}

func onUpdate(old, new interface{}) {
	// Cast the obj as node
	dsOld := old.(*apps.DaemonSet)
	dsNew := new.(*apps.DaemonSet)
	logrus.Infof("Old: %v", dsOld.Labels)
	logrus.Infof("New: %v", dsNew.Labels)
}

func onAdd(obj interface{}) {
	// Cast the obj as node
	ds := obj.(*apps.DaemonSet)
	logrus.Info(ds.Name)

}
