package main

import (
	"github.com/cnvrgctl/pkg/cnvrg"
	v1 "github.com/cnvrgctl/pkg/cnvrg/api/types/v1"
	"testing"
)

func init() {
	setupCommands()
}

func TestCreateCnvrgAppUpgrade(t *testing.T) {
	upgradeSpec := v1.NewCnvrgAppUpgrade("cnvrg/core:3.1.3")
	err := cnvrg.CreateCnvrgAppUpgrade(upgradeSpec)
	if err != nil {
		t.Fatal(err)
	}
	appUpgrade := cnvrg.GetCnvrgAppUpgrade(upgradeSpec.Name)
	if upgradeSpec.Spec.Image != appUpgrade.Spec.Image {
		t.Fatal("images are not equals")
	}
	if upgradeSpec.Spec.Condition != appUpgrade.Spec.Condition {
		t.Fatal("conditions are not equals")
	}
	if upgradeSpec.Spec.CacheImage != appUpgrade.Spec.CacheImage {
		t.Fatal("CacheImage are not equals")
	}
	if upgradeSpec.Spec.CacheDsName != appUpgrade.Spec.CacheDsName {
		t.Fatal("CacheDsName are not equals")
	}
	if upgradeSpec.Spec.CnvrgAppName != appUpgrade.Spec.CnvrgAppName {
		t.Fatal("CacheDsName are not equals")
	}
	t.Cleanup(func() {
		if err := cnvrg.DeleteAllUpgrades(); err != nil {
			t.Fatal(err)
		}
	})
}

func TestCreateAppUpgradeWhileRunningUpgrade(t *testing.T) {
	err := cnvrg.CreateCnvrgAppUpgrade(v1.NewCnvrgAppUpgrade("cnvrg/core:3.1.3"))
	err = cnvrg.CreateCnvrgAppUpgrade(v1.NewCnvrgAppUpgrade("cnvrg/core:3.1.3"))
	if err == nil {
		t.Fatal("Error - shouldn't able to create deploy object")
	}
	t.Cleanup(func() {
		if err := cnvrg.DeleteAllUpgrades(); err != nil {
			t.Fatal(err)
		}
	})
}