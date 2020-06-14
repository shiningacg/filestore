package filestore

import "io"

type API interface {
	Get(uuid string) (File, error)
	Add(reader io.Reader) error
	Remove(uuid string) error
}
