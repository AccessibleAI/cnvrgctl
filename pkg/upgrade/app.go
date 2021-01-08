package upgrade

import (
	"context"
	"fmt"
	"github.com/briandowns/spinner"
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
	"time"
)

const (
	PullImageDsName         = "cnvrg-image-puller"
	CnvrgAppBackupCm        = "cnvrg-app-backup"
	SidekiqDeploymentName   = "sidekiq"
	SearchkiqDeploymentName = "sidekiq-searchkick"
)

type updateEventHandler func(old, new interface{}, stopper chan struct{})

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
	logrus.Infof("gracefully taking down sidekiq")
	cnvrgApp := GetCnvrgApp()
	cnvrgApp.Spec.CnvrgApp.SidekiqReplicas = 0
	cnvrgApp.Spec.CnvrgApp.SidekiqSearchkickReplicas = 0
	UpdateCnvrgApp(cnvrgApp)
	stopper := make(chan struct{})
	go IsSidekiqScaledToZero(stopper)
	<-stopper
	logrus.Infof("sidekiq stopped")
}

func RunApplicationUpgrade() {

	logrus.Infof("running application upgrade")
	cnvrgApp := GetCnvrgApp()
	cnvrgApp.Spec.CnvrgApp.ResourcesRequestEnabled = "false"
	if cnvrgApp.Spec.CnvrgApp.Image == viper.GetString("app-image") {
		logrus.Infof("no need to upgrade, app image already upgraded")
		return
	}
	cnvrgApp.Spec.CnvrgApp.Image = viper.GetString("app-image")
	UpdateCnvrgApp(cnvrgApp)

	stopper := make(chan struct{})
	messages := make(chan string, 100)
	s := spinner.New(spinner.CharSets[27], 50*time.Millisecond)
	cnvrgAppFromBackup := LoadCnvrgAppFromBackup()

	go startSpinner(s, "upgrade is running ", messages)
	go WatchForPods(stopper, "app", []string{"app", "app"}, messages)
	go WatchForDeployments(func(old, new interface{}, stopper chan struct{}) {
		app := new.(*apps.Deployment)
		if app.Name == cnvrgApp.Spec.CnvrgApp.SvcName {
			if app.Spec.Template.Spec.Containers[0].Image == cnvrgApp.Spec.CnvrgApp.Image &&
				cnvrgAppFromBackup.Spec.CnvrgApp.Replicas == app.Spec.Replicas {
				logrus.Info("status: %v:", app.Status)
				s.Stop()
				close(messages)
				close(stopper)
			}
		}
	}, stopper)
	<-stopper
}

func IsSidekiqScaledToZero(stopper chan struct{}) {

	messages := make(chan string, 100)
	s := spinner.New(spinner.CharSets[27], 50*time.Millisecond)
	go startSpinner(s, "scaling sidekiq to 0", messages)
	go WatchForPods(stopper, "app", []string{SidekiqDeploymentName, SearchkiqDeploymentName}, messages)
	for {
		sidekiqDeployment := IsDeploymentScaledTo(
			SidekiqDeploymentName,
			viper.GetString("cnvrg-namespace"),
			0)
		searchkiqDeployment := IsDeploymentScaledTo(
			SearchkiqDeploymentName,
			viper.GetString("cnvrg-namespace"),
			0)

		if sidekiqDeployment && searchkiqDeployment {
			logrus.Debug("%v scaled down to 0", SidekiqDeploymentName)
			logrus.Debug("%v scaled down to 0", SearchkiqDeploymentName)
			s.Stop()
			close(messages)
			close(stopper)
			break
		}
		time.Sleep(5 * time.Second)
	}
}

func IsDeploymentScaledTo(name string, ns string, scaleTo int32) (ready bool) {
	ready = false
	logrus.Debug("getting %v deployment", name)
	_, client := getK8SDefaultClient()
	clientset := client.AppsV1().Deployments(ns)
	d, err := clientset.Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatalf("wasn't able to get %v deploy", name)
	}
	if d.Status.ReadyReplicas == scaleTo {
		logrus.Debug("%v scaled to: %v ", name, scaleTo)
		ready = true
	}
	return
}

func DeployImagePullDaemonSet(cnvrgApp *cnvrgappv1.CnvrgApp, image string) {
	logrus.Debugf("starting image pull daemon set...")
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
	DeleteDaemonSet(PullImageDsName)
	logrus.Debugf("creating pull image DaemonSet: %v", ds)
	_, err := dsClientSet.Create(context.TODO(), ds, metav1.CreateOptions{})
	if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("wasn't able to create image pull ds")
	}
	logrus.Debugf("image pull ds is submitted, waiting to become ready...")
}

func DeleteDaemonSet(name string) {
	_, client := getK8SDefaultClient()
	dsClientSet := client.AppsV1().DaemonSets(viper.GetString("cnvrg-namespace"))
	err := dsClientSet.Delete(context.TODO(), name, metav1.DeleteOptions{})
	if errors.IsNotFound(err) {
		logrus.Debug("nothing to delete, ds %v not found", PullImageDsName)
	} else if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("wasn't able to delete image pull image ds")
	}
}

func startSpinner(s *spinner.Spinner, suffixMessage string, messages <-chan string) {
	s.Suffix = suffixMessage
	s.Color("green")
	s.Start()
	for v := range messages {
		msg := fmt.Sprintf("%v [ %v ]", suffixMessage, v)
		s.Suffix = msg
		s.Restart()
	}
}

func WatchForImagePullDaemonSetReady(ready chan<- bool) {
	_, client := getK8SDefaultClient()
	factory := informers.NewSharedInformerFactory(client, 0)
	informer := factory.Apps().V1().DaemonSets().Informer()
	stopper := make(chan struct{})
	podsWatchStopper := make(chan struct{})
	messages := make(chan string, 100)
	s := spinner.New(spinner.CharSets[27], 50*time.Millisecond)
	go startSpinner(s, "caching app image on all nodes", messages)
	defer runtime.HandleCrash()
	go WatchForPods(podsWatchStopper, "app", []string{PullImageDsName}, messages)
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(old, new interface{}) {
			dsNew := new.(*apps.DaemonSet)
			logrus.Debugf("Status: %v:", dsNew.Status)
			if dsNew.Name == PullImageDsName && dsNew.Status.DesiredNumberScheduled == dsNew.Status.NumberAvailable {
				close(messages)
				close(stopper)
				close(podsWatchStopper)
				s.Stop()
				logrus.Debug("watch DaemonSet %v completed, DS is ready", PullImageDsName)
				DeleteDaemonSet(PullImageDsName)
				ready <- true
			} else {
				logrus.Debug("%v DaemonSet is not ready yet...", PullImageDsName)
			}
		},
	})
	go informer.Run(stopper)
	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		logrus.Debug("timed out waiting for caches to sync or the stopper was closed...")
		return
	}
	<-stopper
}

func WatchForPods(stopper chan struct{}, labelKey string, labelsVal []string, messages chan string) {
	_, client := getK8SDefaultClient()
	factory := informers.NewSharedInformerFactory(client, 0)
	informer := factory.Core().V1().Pods().Informer()
	defer runtime.HandleCrash()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(old, new interface{}) {
			pod := new.(*v1.Pod)
			if val, ok := pod.ObjectMeta.Labels[labelKey]; ok {
				if _, found := Find(labelsVal, val); found == true {
					msg := fmt.Sprintf("pod %v phase: %v", pod.Name, pod.Status.Phase)
					logrus.Debug(msg)
					messages <- msg
					if len(pod.Status.ContainerStatuses) > 0 && pod.Status.ContainerStatuses[0].State.Waiting != nil {
						msg = fmt.Sprintf(
							"pod %v status: %v", pod.Name, pod.Status.ContainerStatuses[0].State.Waiting.Reason)
						logrus.Debug(msg)
						messages <- msg
					}
				}
			}
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			if val, ok := pod.ObjectMeta.Labels[labelKey]; ok {
				if _, found := Find(labelsVal, val); found == true {
					msg := fmt.Sprintf("pod %v deleted", pod.Name)
					logrus.Debug(msg)
					messages <- msg
				}
			}
		},
	})
	go informer.Run(stopper)
	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		logrus.Debug("timed out waiting for caches to sync or the stopper was closed...")
		return
	}
	<-stopper
}

func WatchForDeployments(updateHandler updateEventHandler, stopper chan struct{}) {
	_, client := getK8SDefaultClient()
	factory := informers.NewSharedInformerFactory(client, 0)
	informer := factory.Apps().V1().Deployments().Informer()
	defer runtime.HandleCrash()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(old, new interface{}) { updateHandler(old, new, stopper) },
	})
	go informer.Run(stopper)
	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		logrus.Fatal("timed out waiting for caches to sync or the stopper was closed...")
		return
	}
	<-stopper
}

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
