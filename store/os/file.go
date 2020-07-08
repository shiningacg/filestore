package os

import "os"

type File struct {
	name string
	id   string
	url  string
	*os.File
}

func (f *File) FileName() string {
	return f.name
}

func (f *File) ID() string {
	return f.id
}

func (f *File) Url() string {
	return f.url
}

func fromDBFile(file *DBFile) (*File, error) {
	f, err := os.Open(file.Path)
	if err != nil {
		return nil, err
	}
	return &File{
		name: file.Name,
		id:   file.UUID,
		File: f,
	}, nil
}
