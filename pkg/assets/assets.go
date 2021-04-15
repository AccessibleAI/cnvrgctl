package assets

import "github.com/markbates/pkger"

func informPkger() {
	pkger.Include("/pkg/assets/cluster-setup.sh")
	pkger.Include("/pkg/assets/cluster.tpl")
}
