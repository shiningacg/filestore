module github.com/shiningacg/filestore

go 1.15

replace (
	github.com/shiningacg/Services => /Users/shlande/go/src/github.com/shiningacg/Services
	go.etcd.io/etcd => go.etcd.io/etcd v0.0.0-20200824191128-ae9734ed278b
	google.golang.org/grpc => google.golang.org/grpc v1.29.1
)

require (
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/boltdb/bolt v1.3.1
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/go-redis/redis/v8 v8.2.3
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.2
	github.com/shiningacg/Services v0.0.0-00010101000000-000000000000
	github.com/shiningacg/mygin v0.0.0-20201003081440-00529e907d03
	github.com/shiningacg/mygin-frame-libs v0.0.0-20200801133652-d3ee76596824
	github.com/shiningacg/sn-ipfs v0.0.0-20200924124624-1bb5619e1f1a
	github.com/shirou/gopsutil v2.20.9+incompatible
	go.etcd.io/etcd v3.3.25+incompatible
	google.golang.org/grpc v1.32.0
)
