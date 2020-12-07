package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CnvrgAppInterface interface {
	List(opts metav1.ListOptions) (*CnvrgAppList, error)
}
