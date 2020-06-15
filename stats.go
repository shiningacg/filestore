package filestore

type Stats interface {
	Space() Space
	Bandwidth() Bandwidth
}

type Space struct {
	Total uint64
	Free  uint64
	Used  uint64
}

type Bandwidth struct {
	Upload        uint64
	Download      uint64
	TotalUpload   uint64
	TotalDownload uint64
}
