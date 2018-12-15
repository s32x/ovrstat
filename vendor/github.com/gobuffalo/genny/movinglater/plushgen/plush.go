package plushgen

import (
	"github.com/gobuffalo/plushgen"
	"github.com/markbates/oncer"
)

func init() {
	oncer.Deprecate(0, "github.com/gobuffalo/genny/movinglater/plushgen", "Use github.com/gobuffalo/plushgen instead.")
}

// Transformer will plushify any file that has a ".plush" extension
var Transformer = plushgen.Transformer
