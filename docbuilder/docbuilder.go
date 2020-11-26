package docbuilder

import "context"

type DocBuilder interface {
	Build(ctx context.Context, srcpath, distpath string, force bool) (err error)
}
