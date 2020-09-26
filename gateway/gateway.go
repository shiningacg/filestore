package gateway

import (
	"context"
	fs "github.com/shiningacg/filestore"
)

type Gateway interface {
	Run(ctx context.Context) error
	BandWidth() *fs.Bandwidth
	SetStore(store fs.FileStore)
	Host() string
}
