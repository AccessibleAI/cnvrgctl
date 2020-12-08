package k8s

import (
	"context"
	cnvrgappv1 "github.com/cnvrgctl/pkg/cnvrgapp/api/types/v1"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cnvrgappV1client "github.com/cnvrgctl/pkg/cnvrgapp/clientset/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
)

func GetNodes() *v1.NodeList {
	config, err := clientcmd.BuildConfigFromFlags("", viper.GetString("kubeconfig"))
	if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("Error building kubeconfig")
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("Error creating client")
	}
	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Errorf("Can't fetch K8S cluster nodes")
		logrus.Debug(err.Error())
	}
	return nodes
}

// https://www.martin-helmich.de/en/blog/kubernetes-crd-client.html
func GetCnvrgApp() {
	config, err := clientcmd.BuildConfigFromFlags("", viper.GetString("kubeconfig"))
	if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("Error building kubeconfig")
	}
	if err := cnvrgappv1.AddToScheme(scheme.Scheme); err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("Error registering cnvrgapp CR")

	}
	clientSet, err := cnvrgappV1client.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	//cnvrgapps, err := clientSet.CnvrgApps("cnvrg").List(context.TODO(), metav1.ListOptions{})
	cnvrgapp, err := clientSet.CnvrgApps("cnvrg").Get(context.TODO(), "cnvrg-app", metav1.GetOptions{})
	cnvrgapp.Spec.CnvrgApp.Image = "cnvrg/app:master-3553"
	result, err := clientSet.CnvrgApps("cnvrg").Update(context.TODO(), cnvrgapp, metav1.UpdateOptions{})

	logrus.Info(cnvrgapp)
	logrus.Info(result)
	//config, err := clientcmd.BuildConfigFromFlags("", viper.GetString("kubeconfig"))
	//if err != nil {
	//	logrus.Debug(err.Error())
	//	logrus.Fatal("Error building kubeconfig")
	//}
	////client, err := kubernetes.NewForConfig(config)
	////if err != nil {
	////	logrus.Debug(err.Error())
	////	logrus.Fatal("Error creating client")
	////}
	//if err := cnvrgappv1.AddToScheme(scheme.Scheme); err != nil {
	//	logrus.Debug(err.Error())
	//	logrus.Fatal("Error registering cnvrgapp CR")
	//
	//}
	//
	//crdConfig := *config
	//crdConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: cnvrgappv1.GroupName, Version: cnvrgappv1.GroupVersion}
	//crdConfig.APIPath = "/apis"
	//crdConfig.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	//crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()
	//exampleRestClient, err := rest.UnversionedRESTClientFor(&crdConfig)
	//result := cnvrgappv1.CnvrgAppList{}
	//if err := exampleRestClient.Get().Resource("cnvrgapps").Do(context.TODO()).Into(&result); err != nil {
	//	logrus.Debug(err.Error())
	//	logrus.Fatal("error fetch cnvrgapp CR")
	//}
	//logrus.Info(result)
}
