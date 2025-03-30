package tg

import (
	"log"
	"os"

	"github.com/google/wire"
	"github.com/thekhanj/csdmpro/tg/controllers"
	"github.com/thekhanj/tgool"
)

type TgControllers []tgool.Controller

func ProvideControllers() TgControllers {
	return TgControllers{
		&controllers.WatchlistController{},
	}
}

func ProvideTg(controllers TgControllers) *Server {
	serverBuilder := ServerBuilder{}

	serverBuilder.
		WithToken(os.Getenv("API_TOKEN")).
		WithControllers(controllers...)

	socks_proxy := os.Getenv("http_proxy")
	if socks_proxy != "" {
		serverBuilder.WithProxy(socks_proxy)
	}

	s, err := serverBuilder.Build()
	if err != nil {
		log.Fatal(err)
	}

	return s
}

var TgModule = wire.NewSet(ProvideTg, ProvideControllers)
