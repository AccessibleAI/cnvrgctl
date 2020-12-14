package upgrade

import (
	"context"
	cnvrgappv1 "github.com/cnvrgctl/pkg/cnvrgapp/api/types/v1"
	cnvrgappV1client "github.com/cnvrgctl/pkg/cnvrgapp/clientset/v1"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

const (
	PullImageDsName         = "cnvrg-image-puller"
	CnvrgAppBackupCm        = "cnvrg-app-backup"
	SidekiqDeploymentName   = "sidekiq"
	SearchkiqDeploymentName = "sidekiq-searchkick"
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

func GetNodesMetrics() {
	config, _ := getK8SDefaultClient()
	mc, err := metrics.NewForConfig(config)
	if err != nil {
		logrus.Errorf("can't fetch K8S cluster nodes metrics")
		logrus.Debug(err.Error())
	}
	nodes, err := mc.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{})
	for _, node := range nodes.Items {
		logrus.Infof("node: %v cpu: %v, memory: %v",
			node.Name,
			node.Usage.Cpu().AsDec(),
			node.Usage.Memory().AsDec())
		logrus.Info(node.Name)
	}
}

func GetNodes() *v1.NodeList {
	_, client := getK8SDefaultClient()
	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Errorf("can't fetch K8S cluster nodes")
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
	return
}

func GetAppDeployment(webAppDeploymentName string) (appDeploy *apps.Deployment) {
	logrus.Info("getting webapp deployment")
	_, client := getK8SDefaultClient()
	clientset := client.AppsV1().Deployments(viper.GetString("cnvrg-namespace"))
	appDeploy, err := clientset.Get(context.TODO(), webAppDeploymentName, metav1.GetOptions{})
	if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("error fetching app deployment")
	}
	return
}

func LoadCnvrgAppFromBackup() *cnvrgappv1.CnvrgApp {
	cmNs := "default"
	_, clientset := getK8SDefaultClient()
	backupCm, err := clientset.CoreV1().ConfigMaps(cmNs).Get(context.TODO(), CnvrgAppBackupCm, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		logrus.Warnf("cnvrgapp backup not found..")
	} else if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("error getting cnvrgapp backup configmap")
	}
	var cnvrgApp cnvrgappv1.CnvrgApp
	if data, ok := backupCm.Data["cnvrgApp"]; ok {
		if err := json.Unmarshal([]byte(data), &cnvrgApp); err != nil {
			logrus.Debug(err.Error())
			logrus.Fatal("can't unmarshal backup cnvrgapp")
		}
	} else {
		logrus.Warnf("cnvrgapp key not found in configmap: %v", CnvrgAppBackupCm)
	}
	return &cnvrgApp
}

func UpdateCnvrgApp(cnvrgApp *cnvrgappv1.CnvrgApp) {
	if cnvrgApp.Name == "" {
		logrus.Debug(cnvrgApp)
		logrus.Fatal("can't update empty cnvrgapp")
	}
	config, _ := getK8SDefaultClient()
	if err := cnvrgappv1.AddToScheme(scheme.Scheme); err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("Error registering cnvrgapp CR")
	}
	clientSet, err := cnvrgappV1client.NewForConfig(config)
	if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("error creating cnvrgappv1 clientset")
	}
	ns := viper.GetString("cnvrg-namespace")
	cnvrgAppForUpdate := GetCnvrgApp()
	cnvrgAppForUpdate.Spec = cnvrgApp.Spec
	_, err = clientSet.CnvrgApps(ns).Update(context.TODO(), cnvrgAppForUpdate, metav1.UpdateOptions{})
	if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("error updating cnvrgapp")
	}
}

func BackupCnvrgApp() {
	logrus.Infof("creating backup cnvrgapp configmap: %v", CnvrgAppBackupCm)
	cmNs := "default"
	cnvrgApp := GetCnvrgApp()
	cnvrgAppForBackup := cnvrgappv1.CnvrgApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cnvrgApp.ObjectMeta.Name,
			Namespace: cnvrgApp.ObjectMeta.Namespace,
		},
		Spec: cnvrgApp.Spec,
	}
	b, err := json.Marshal(cnvrgAppForBackup)
	if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("error marshaling app spec")
	}
	_, clientset := getK8SDefaultClient()
	cmData := map[string]string{"cnvrgApp": string(b)}
	cm := &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      CnvrgAppBackupCm,
			Namespace: "default",
		},
		Data: cmData,
	}
	_, err = clientset.CoreV1().ConfigMaps(cmNs).Create(context.TODO(), cm, metav1.CreateOptions{})
	if errors.IsAlreadyExists(err) {
		logrus.Warnf("cnvrgapp backup cm already exists, not going to create new one until current upgrade will be finished ")
	} else if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("wasn't able to create backup cm")
	}
}

func SidekiqGracefulShutdown() {
	cnvrgApp := GetCnvrgApp()
	cnvrgApp.Spec.CnvrgApp.SidekiqReplicas = 0
	cnvrgApp.Spec.CnvrgApp.SidekiqSearchkickReplicas = 0
	UpdateCnvrgApp(cnvrgApp)
	stopper := make(chan bool)
	go WatchForDeploymentScaleToZero(stopper)
	<-stopper
}

func WatchForDeploymentScaleToZero(stopper chan<- bool) {
	logrus.Info("getting sidekiq deployment")
	_, client := getK8SDefaultClient()
	clientset := client.AppsV1().Deployments(viper.GetString("cnvrg-namespace"))
	sidekiqDeploy, err := clientset.Get(context.TODO(), SidekiqDeploymentName, metav1.GetOptions{})
	if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("wasn't able to get sidekiq deploy")
	}
	searchkiqDeploy, err := clientset.Get(context.TODO(), SearchkiqDeploymentName, metav1.GetOptions{})
	if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("wasn't able to get sidekiq-searchkick deploy")
	}
	if searchkiqDeploy.Status.ReadyReplicas == 0 && sidekiqDeploy.Status.ReadyReplicas == 0 {
		logrus.Info("sidekiq deployment is scaled down")
		logrus.Info("sidekiq-searchkic deployment is scaled down")
		stopper <- true
	}

	if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("error fetching app deployment")
	}
	return
}

func DeployImagePullDaemonSet(cnvrgApp *cnvrgappv1.CnvrgApp, image string) {
	logrus.Info("starting image pull daemon set...")
	specSelectorLabels := map[string]string{"app": PullImageDsName}
	command := []string{"/bin/bash", "-c", "sleep inf"}
	appDeployment := GetAppDeployment(cnvrgApp.Spec.CnvrgApp.SvcName)
	logrus.Debugf("image cache ds using pull secret: %v", appDeployment.Spec.Template.Spec.ImagePullSecrets)
	logrus.Debugf("image cache ds using toleration: %v", appDeployment.Spec.Template.Spec.Tolerations)
	logrus.Debugf("image cache ds using node selector: %v", appDeployment.Spec.Template.Spec.NodeSelector)
	ds := &apps.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: viper.GetString("cnvrg-namespace"),
			Name:      PullImageDsName,
		},
		Spec: apps.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: specSelectorLabels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: specSelectorLabels,
				},
				Spec: v1.PodSpec{
					ImagePullSecrets: appDeployment.Spec.Template.Spec.ImagePullSecrets,
					Tolerations:      appDeployment.Spec.Template.Spec.Tolerations,
					NodeSelector:     appDeployment.Spec.Template.Spec.NodeSelector,
					Containers: []v1.Container{
						{
							Name:    PullImageDsName,
							Image:   image,
							Command: command,
						},
					},
				},
			},
		},
	}
	_, client := getK8SDefaultClient()
	dsClientSet := client.AppsV1().DaemonSets(viper.GetString("cnvrg-namespace"))
	err := dsClientSet.Delete(context.TODO(), PullImageDsName, metav1.DeleteOptions{})
	if errors.IsNotFound(err) {
		logrus.Warnf("nothing to delete, ds %v not found", PullImageDsName)
	} else if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("wasn't able to delete image pull image ds")
	}
	logrus.Debugf("creating pull image DaemonSet: %v", ds)
	_, err = dsClientSet.Create(context.TODO(), ds, metav1.CreateOptions{})
	if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("wasn't able to create image pull ds")
	}

	logrus.Info("image pull ds is submitted, waiting to become ready...")
}

func WatchForImagePullDaemonSetReady(ready chan<- bool) {
	_, client := getK8SDefaultClient()
	factory := informers.NewSharedInformerFactory(client, 0)
	informer := factory.Apps().V1().DaemonSets().Informer()
	stopper := make(chan struct{})
	podsWatchStopper := make(chan struct{})
	defer runtime.HandleCrash()
	go WatchForPodReady(podsWatchStopper, "app", PullImageDsName)
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(old, new interface{}) {
			dsNew := new.(*apps.DaemonSet)
			logrus.Debugf("Status: %v:", dsNew.Status)
			if dsNew.Name == PullImageDsName && dsNew.Status.DesiredNumberScheduled == dsNew.Status.NumberAvailable {
				logrus.Infof("watch DaemonSet %v completed, DS is ready", PullImageDsName)
				close(stopper)
				close(podsWatchStopper)
				ready <- true
			} else {
				logrus.Infof("%v DaemonSet is not ready yet...", PullImageDsName)
			}
		},
	})
	go informer.Run(stopper)
	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		logrus.Fatal("timed out waiting for caches to sync")
		close(stopper)
		return
	}
	<-stopper
}

func WatchForPodReady(stopper chan struct{}, labelKey string, labelVal string) {
	_, client := getK8SDefaultClient()
	factory := informers.NewSharedInformerFactory(client, 0)
	informer := factory.Core().V1().Pods().Informer()
	defer runtime.HandleCrash()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(old, new interface{}) {
			pod := new.(*v1.Pod)
			if val, ok := pod.ObjectMeta.Labels[labelKey]; ok {
				if val == labelVal {
					logrus.Infof("pod %v phase: %v", pod.Name, pod.Status.Phase)
					if len(pod.Status.ContainerStatuses) > 0 && pod.Status.ContainerStatuses[0].State.Waiting != nil {
						logrus.Infof("pod %v status: %v", pod.Name,
							pod.Status.ContainerStatuses[0].State.Waiting.Reason)
					}
				}
			}
		},
	})
	go informer.Run(stopper)
	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		logrus.Fatal("timed out waiting for caches to sync")
		close(stopper)
		return
	}
	<-stopper
}
