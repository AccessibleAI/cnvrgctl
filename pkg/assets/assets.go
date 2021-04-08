package assets

import "github.com/markbates/pkger"

func informPkger() {
	pkger.Include("/pkg/assets/rke")
	pkger.Include("/pkg/assets/k9s")
	pkger.Include("/pkg/assets/kubectl")
}
