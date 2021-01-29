package v1

import (
	cnvrgappv1 "github.com/cnvrgctl/pkg/cnvrgapp/api/types/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

// https://www.martin-helmich.de/en/blog/kubernetes-crd-client.html
type CnvrgAppV1Interface interface {
	CnvrgApps(namespace string) CnvrgAppInterface
}

type CnvrgAppUpgradeV1Interface interface {
	CnvrgAppUpgrades(namespace string) CnvrgAppUpgradeInterface
}

type CnvrgAppV1Client struct {
	restClient rest.Interface
}

type CnvrgAppUpgradeV1Client struct {
	restClient rest.Interface
}

func NewForConfigCnvrgApp(c *rest.Config) (*CnvrgAppV1Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: cnvrgappv1.GroupName, Version: cnvrgappv1.GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &CnvrgAppV1Client{restClient: client}, nil
}

func NewForConfigCnvrgAppUpgrade(c *rest.Config) (*CnvrgAppUpgradeV1Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: cnvrgappv1.GroupName, Version: cnvrgappv1.GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &CnvrgAppUpgradeV1Client{restClient: client}, nil
}

func (c *CnvrgAppV1Client) CnvrgApps(namespace string) CnvrgAppInterface {
	return &cnvrgappClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

func (c *CnvrgAppUpgradeV1Client) CnvrgAppUpgrades(namespace string) CnvrgAppUpgradeInterface {
	return &cnvrgappupgradeClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
