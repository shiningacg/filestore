package filestore

import (
	"context"
)

type Gateway interface {
	Run(ctx context.Context) error
	BandWidth() *Bandwidth
	SetStore(store FileFS)
	Host() string
}
