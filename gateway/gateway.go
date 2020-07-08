package gateway

import (
	store "github.com/shiningacg/filestore"
	"net/http"
)

type Gateway struct {
	store.API
}

func (g *Gateway) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	panic("implement me")
}

func (g *Gateway) Gateway() store.Gateway {
	panic("implement me")
}

func (g *Gateway) Run(addr string) error {
	return http.ListenAndServe(addr, g)
}
