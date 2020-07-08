package gateway

import (
	"fmt"
	store "github.com/shiningacg/filestore"
	"net/http"
)

func NewGateway(addr string, api store.API) *Gateway {
	return &Gateway{
		addr: addr,
		api:  api,
	}
}

type Gateway struct {
	addr string
	api  store.API
}

func (g *Gateway) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	panic("implement me")
}

func (g *Gateway) Gateway() store.Gateway {
	panic("implement me")
}

func (g *Gateway) Run() error {
	return http.ListenAndServe(g.addr, g)
}

func GetBaseUrl(addr string) string {
	return fmt.Sprintf("http://%v/get/", addr)
}
