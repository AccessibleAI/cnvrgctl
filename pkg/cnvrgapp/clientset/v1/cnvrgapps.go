package v1

import (
	"context"
	cnvrgappv1 "github.com/cnvrgctl/pkg/cnvrgapp/api/types/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type CnvrgAppInterface interface {
	List(ctx context.Context, opts metav1.ListOptions) (*cnvrgappv1.CnvrgAppList, error)
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*cnvrgappv1.CnvrgApp, error)
	Update(ctx context.Context, cnvrgapp *cnvrgappv1.CnvrgApp, opts metav1.UpdateOptions) (*cnvrgappv1.CnvrgApp, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
}

type cnvrgappClient struct {
	restClient rest.Interface
	ns         string
}

func (c *cnvrgappClient) List(ctx context.Context, opts metav1.ListOptions) (result *cnvrgappv1.CnvrgAppList, err error) {
	result = &cnvrgappv1.CnvrgAppList{}
	err = c.restClient.
		Get().
		Namespace(c.ns).
		Resource("cnvrgapps").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

func (c *cnvrgappClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (result *cnvrgappv1.CnvrgApp, err error) {
	result = &cnvrgappv1.CnvrgApp{}
	err = c.restClient.
		Get().
		Namespace(c.ns).
		Resource("cnvrgapps").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

func (c *cnvrgappClient) Update(ctx context.Context, cnvrgapp *cnvrgappv1.CnvrgApp, opts metav1.UpdateOptions) (result *cnvrgappv1.CnvrgApp, err error) {
	result = &cnvrgappv1.CnvrgApp{}
	err = c.restClient.
		Put().
		Namespace(c.ns).
		Resource("cnvrgapps").
		Name(cnvrgapp.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(cnvrgapp).
		Do(ctx).
		Into(result)
	return
}

func (c *cnvrgappClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("cnvrgapps").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}
