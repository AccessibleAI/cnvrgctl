package assets

import "github.com/markbates/pkger"

func informPkger() {
	pkger.Include("/pkg/assets/")
	pkger.Include("/pkg/assets/k9s")
	pkger.Include("/pkg/assets/kubectl")
	pkger.Include("/pkg/assets/cluster.tpl")
	pkger.Include("/pkg/assets/cnvrg-crds.yaml")
	pkger.Include("/pkg/assets/cnvrg-operator.yaml")
}
