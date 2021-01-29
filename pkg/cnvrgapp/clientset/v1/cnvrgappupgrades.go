package v1

import (
	"context"
	cnvrgappv1 "github.com/cnvrgctl/pkg/cnvrgapp/api/types/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type CnvrgAppUpgradeInterface interface {
	List(ctx context.Context, opts metav1.ListOptions) (*cnvrgappv1.CnvrgAppUpgradeList, error)
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*cnvrgappv1.CnvrgAppUpgrade, error)
	Update(ctx context.Context, cnvrgapp *cnvrgappv1.CnvrgAppUpgrade, opts metav1.UpdateOptions) (*cnvrgappv1.CnvrgAppUpgrade, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
}


type cnvrgappupgradeClient struct {
	restClient rest.Interface
	ns         string
}

func (c *cnvrgappupgradeClient) List(ctx context.Context, opts metav1.ListOptions) (result *cnvrgappv1.CnvrgAppUpgradeList, err error) {
	result = &cnvrgappv1.CnvrgAppUpgradeList{}
	err = c.restClient.
		Get().
		Namespace(c.ns).
		Resource("cnvrgappupgrades").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

func (c *cnvrgappupgradeClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (result *cnvrgappv1.CnvrgAppUpgrade, err error) {
	result = &cnvrgappv1.CnvrgAppUpgrade{}
	err = c.restClient.
		Get().
		Namespace(c.ns).
		Resource("cnvrgappupgrades").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

func (c *cnvrgappupgradeClient) Update(ctx context.Context, cnvrgapp *cnvrgappv1.CnvrgAppUpgrade, opts metav1.UpdateOptions) (result *cnvrgappv1.CnvrgAppUpgrade, err error) {
	result = &cnvrgappv1.CnvrgAppUpgrade{}
	err = c.restClient.
		Put().
		Namespace(c.ns).
		Resource("cnvrgappupgrades").
		Name(cnvrgapp.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(cnvrgapp).
		Do(ctx).
		Into(result)
	return
}

func (c *cnvrgappupgradeClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("cnvrgappupgrades").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}
