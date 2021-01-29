package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CnvrgAppUpgradeSpec struct {
	Condition    string `json:"condition"`
	CacheDsName  string `json:"cacheDsName"`
	CnvrgAppName string `json:"cnvrgAppName"`
	Image        string `json:"image"`
	CacheImage   string `json:"cacheImage"`
}

type CnvrgAppUpgradeStatus struct {
	Status struct {
		Status string `json:"status"`
	} `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type CnvrgAppUpgrade struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CnvrgAppUpgradeSpec   `json:"spec,omitempty"`
	Status CnvrgAppUpgradeStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type CnvrgAppUpgradeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CnvrgAppUpgrade `json:"items"`
}
