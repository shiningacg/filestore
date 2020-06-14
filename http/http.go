package http

import (
	filestore "filesys"
	"fmt"
	"github.com/shiningacg/apicore"
)

/*
	http功能，主要负责流控和一些简单的记录
*/

var Store filestore.Store

type Recode struct {
	RequestID string
	Ip        string
	FileID    string
	Bandwidth uint64
	Finish    bool
}

var InputChan chan<- Recode

func SetStore(store filestore.Store) {
	Store = store
}

func SetRecodChan(input chan<- Recode) {
	InputChan = input
}

func RunGateway(port int) {
	if Store == nil {
		panic("no store find!")
	}
	go func() {
		apicore.Run(fmt.Sprintf(":%v", port))
	}()
}
