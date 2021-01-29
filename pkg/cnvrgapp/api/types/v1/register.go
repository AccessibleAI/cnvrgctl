package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// +kubebuilder:object:generate=true
// +groupName=mlops.cnvrg.io
const GroupName = "mlops.cnvrg.io"
const GroupVersion = "v1"

var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: GroupVersion}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion, &CnvrgApp{}, &CnvrgAppList{})
	scheme.AddKnownTypes(SchemeGroupVersion, &CnvrgAppUpgrade{}, &CnvrgAppUpgradeList{})
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
