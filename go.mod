module github.com/shiningacg/filestore

go 1.15

replace (
	github.com/coreos/bbolt v1.3.4 => go.etcd.io/bbolt v1.3.4
	github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
	google.golang.org/grpc v1.31.1 => google.golang.org/grpc v1.26.0
)

replace github.com/shiningacg/sn-ipfs => /Users/shlande/go/src/github.com/shiningacg/ipfs-filestore

replace github.com/shiningacg/kit-file => /Users/shlande/go/src/github.com/shlande/kit-file

replace github.com/shiningacg/filestore => /Users/shlande/go/src/filestore

require (
	github.com/boltdb/bolt v1.3.1
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/go-redis/redis/v8 v8.0.0-beta.10
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.2
	github.com/shiningacg/apicore v0.0.0-20200614083315-95323a77c7ee
	github.com/shiningacg/kit-file v0.0.0-00010101000000-000000000000
	github.com/shiningacg/mygin v0.0.0-20200716042016-83c605845398
	github.com/shiningacg/mygin-frame-libs v0.0.0-20200801133652-d3ee76596824
	github.com/shiningacg/sn-ipfs v0.0.0-00010101000000-000000000000
	github.com/shirou/gopsutil v2.20.8+incompatible
	go.etcd.io/etcd v3.3.25+incompatible
	google.golang.org/grpc v1.31.1
)
