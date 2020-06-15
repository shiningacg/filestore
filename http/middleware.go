package http

import (
	"github.com/shiningacg/apicore"
	store "github.com/shiningacg/filestore"
)

const MDN = "GV"

func NewMiddleware(store store.Store, input chan<- *Record) *Global {
	return &Global{store: store, input: input}
}

type Global struct {
	store store.Store
	input chan<- *Record
}

func (m *Global) Before(ctx apicore.Context) {
	ctx.SetValue(MDN, map[string]interface{}{"store": m.store, "input": m.input})
}

func (m *Global) After(ctx apicore.Context) {
	return
}

func (m *Global) Index() int {
	return 23
}

func Store(ctx apicore.Conn) store.Store {
	return ctx.(apicore.Context).Value(MDN).(map[string]interface{})["store"].(store.Store)
}

func Chan(ctx apicore.Conn) chan<- *Record {
	return ctx.(apicore.Context).Value(MDN).(map[string]interface{})["input"].(chan<- *Record)
}
