package v1

import (
	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
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

func NewCnvrgAppUpgrade(namespace string, cnvrgAppName string, image string, cacheImage string) *CnvrgAppUpgrade {

	cnvrgAppUpgrade := CnvrgAppUpgrade{
		TypeMeta:   metav1.TypeMeta{APIVersion: "mlops.cnvrg.io/v1", Kind: "CnvrgAppUpgrade"},
		ObjectMeta: metav1.ObjectMeta{Name: getUpgradeSpecName(image), Namespace: namespace},
		Spec: CnvrgAppUpgradeSpec{
			Condition:    "upgrade",
			CacheDsName:  "app-image-cache",
			CnvrgAppName: cnvrgAppName,
			Image:        image,
			CacheImage:   cacheImage,
		},
	}
	return &cnvrgAppUpgrade
}

func getUpgradeSpecName(image string) string {
	specName := "upgrade-" + uuid.New().String()
	specNameStrArr := strings.Split(image, ":")
	if len(specNameStrArr) == 2 {
		specName = specNameStrArr[1]
	}
	return specName
}

