package mock

import "os"

type File struct {
	url  string
	id   string
	name string
	*os.File
}

func (f File) FileName() string {
	return f.name
}

func (f File) ID() string {
	return f.id
}

func (f File) Url() string {
	return "http://127.0.0.1:8080/get/aaaa"
}
