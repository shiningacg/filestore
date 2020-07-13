package remote

type File struct {
	uuid string
	name string
	size uint64
	url  string
}

// remote无法直接进行读取操作
func (f File) Read(p []byte) (n int, err error) {
	panic("implement me")
}

// remote无法进行读取操作
func (f File) Seek(offset int64, whence int) (int64, error) {
	panic("implement me")
}

func (f File) Close() error {
	return nil
}

func (f File) FileName() string {
	return f.name
}

func (f File) ID() string {
	return f.uuid
}

func (f File) Url() string {
	return f.url
}
