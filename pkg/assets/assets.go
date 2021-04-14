package assets

import "github.com/markbates/pkger"

func informPkger() {
	pkger.Include("/pkg/assets/cluster-setup.sh")
	pkger.Include("/pkg/assets/cluster.tpl")
	pkger.Include("/pkg/assets/cnvrg-crds.yaml")
	pkger.Include("/pkg/assets/cnvrg-operator.yaml")
}
