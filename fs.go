package filestore

type InfoFS interface {
	Get(uuid string) (BaseFile, error)
	Add(file BaseFile) error
	Remove(uuid string) error
}

type FileFS interface {
	Get(uuid string) (ReadableFile, error)
	Add(file ReadableFile) error
	Remove(uuid string) error
}
