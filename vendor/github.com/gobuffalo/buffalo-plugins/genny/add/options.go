package add

import (
	"os"

	"github.com/gobuffalo/buffalo-plugins/plugins/plugdeps"
	"github.com/gobuffalo/meta"
	"github.com/pkg/errors"
)

type Options struct {
	App     meta.App
	Plugins []plugdeps.Plugin
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	if opts.App.IsZero() {
		pwd, err := os.Getwd()
		if err != nil {
			return errors.WithStack(err)
		}
		opts.App = meta.New(pwd)
	}
	if len(opts.Plugins) == 0 {
		plugs, err := plugdeps.List(opts.App)
		if err != nil && (errors.Cause(err) != plugdeps.ErrMissingConfig) {
			return errors.WithStack(err)
		}
		opts.Plugins = plugs.List()
	}
	return nil
}
