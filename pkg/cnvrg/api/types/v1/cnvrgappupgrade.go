package v1

import (
	"github.com/google/uuid"
	"github.com/spf13/viper"
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
	Status      string      `json:"status"`
	CnvrgBackup interface{} `json:"cnvrgBackup"`
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

func NewCnvrgAppUpgrade(image string) *CnvrgAppUpgrade {

	cnvrgAppUpgrade := CnvrgAppUpgrade{
		TypeMeta: metav1.TypeMeta{APIVersion: "mlops.cnvrg.io/v1", Kind: "CnvrgAppUpgrade"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      getUpgradeSpecName(image),
			Namespace: viper.GetString("cnvrg-namespace"),
		},
		Spec: CnvrgAppUpgradeSpec{
			Condition:    viper.GetString("condition"),
			CacheDsName:  viper.GetString("cacheDsName"),
			CnvrgAppName: viper.GetString("cnvrgAppName"),
			Image:        image,
			CacheImage:   viper.GetString("cacheImage"),
		},
		Status: CnvrgAppUpgradeStatus{
			Status: "initiating",
		},
	}
	return &cnvrgAppUpgrade
}

func getUpgradeSpecName(image string) string {
	if viper.GetString("upgrade-name") != "" {
		return viper.GetString("upgrade-name")
	}
	specName := "upgrade-"
	specNameStrArr := strings.Split(image, ":")
	if len(specNameStrArr) == 2 {
		specName = specNameStrArr[1]
	}
	return specName + "-" + strings.Split(uuid.New().String(), "-")[0]
}
